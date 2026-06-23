using MongoDB.Bson;
using MongoDB.Driver;
using ServerCSharp.Models;
using ServerCSharp.Services;
using ServerCSharp.Utils;

namespace ServerCSharp.Endpoints;

public static class UserEndpoints
{
    public static void MapUserEndpoints(this IEndpointRouteBuilder app)
    {
        app.MapGet("/api/users", GetRecommendedUsers);
        app.MapGet("/api/users/me", GetMeAuth);
        app.MapGet("/api/users/friends", GetMyFriends);
        app.MapPost("/api/users/friend-request/{id}", SendFriendRequest);
        app.MapPut("/api/users/friend-request/{id}/accept", AcceptFriendRequest);
        app.MapGet("/api/users/friend-requests", GetFriendRequests);
        app.MapGet("/api/users/outgoing-friend-requests", GetOutgoingFriendRequest);
    }

    private static async Task<IResult> GetRecommendedUsers(HttpContext context, DatabaseService db)
    {
        if (context.Items["userId"] is not string userId)
        {
            return Results.Unauthorized();
        }

        if (!ObjectId.TryParse(userId, out var objectId))
        {
            return Results.BadRequest("Invalid id");
        }

        var me = await db.Users.Find(Builders<User>.Filter.Eq("_id", objectId)).FirstOrDefaultAsync();
        if (me == null)
        {
            return Results.NotFound("No user found");
        }

        var filters = new List<FilterDefinition<User>>
        {
            Builders<User>.Filter.Ne("_id", objectId),
            Builders<User>.Filter.Eq("isonboarded", true)
        };

        if (me.Friends.Count > 0)
        {
            var friendObjectIds = me.Friends
                .Where(id => ObjectId.TryParse(id, out _))
                .Select(id => new ObjectId(id))
                .ToList();

            if (friendObjectIds.Count > 0)
            {
                filters.Add(Builders<User>.Filter.Nin("_id", friendObjectIds));
            }
        }

        var combinedFilter = Builders<User>.Filter.And(filters);
        var users = await db.Users.Find(combinedFilter).ToListAsync();

        return Results.Ok(users);
    }

    private static async Task<IResult> GetMeAuth(HttpContext context, DatabaseService db)
    {
        if (context.Items["userId"] is not string userId)
        {
            return Results.Unauthorized();
        }

        if (!ObjectId.TryParse(userId, out var objectId))
        {
            return Results.NotFound("No user found");
        }

        var user = await db.Users.Find(Builders<User>.Filter.Eq("_id", objectId)).FirstOrDefaultAsync();
        if (user == null)
        {
            return Results.NotFound("No user found");
        }

        return Results.Ok(user);
    }

    private static async Task<IResult> GetMyFriends(HttpContext context, DatabaseService db)
    {
        if (context.Items["userId"] is not string userId)
        {
            return Results.Unauthorized();
        }

        if (!ObjectId.TryParse(userId, out var objectId))
        {
            return Results.BadRequest("Invalid id");
        }

        var me = await db.Users.Find(Builders<User>.Filter.Eq("_id", objectId)).FirstOrDefaultAsync();
        if (me == null)
        {
            return Results.NotFound("No user found");
        }

        if (me.Friends.Count == 0)
        {
            return Results.Ok(new List<User>());
        }

        var friendObjectIds = me.Friends
            .Where(id => ObjectId.TryParse(id, out _))
            .Select(id => new ObjectId(id))
            .ToList();

        if (friendObjectIds.Count == 0)
        {
            return Results.Ok(new List<User>());
        }

        var friends = await db.Users.Find(Builders<User>.Filter.In("_id", friendObjectIds)).ToListAsync();

        return Results.Ok(friends);
    }

    private static async Task<IResult> SendFriendRequest(
        string id,
        HttpContext context,
        DatabaseService db)
    {
        if (context.Items["userId"] is not string userId)
        {
            return Results.Unauthorized();
        }

        if (string.IsNullOrEmpty(id))
        {
            return Results.BadRequest("No recipientId provied");
        }

        if (userId == id)
        {
            return Results.BadRequest("You can't send friend request to yourself");
        }

        if (!ObjectId.TryParse(userId, out var senderObjId) || !ObjectId.TryParse(id, out var recipientObjId))
        {
            return Results.BadRequest("Invalid id");
        }

        // Check if already friends
        var areFriends = await db.Users
            .Find(Builders<User>.Filter.And(
                Builders<User>.Filter.Eq("_id", recipientObjId),
                Builders<User>.Filter.AnyEq("friends", userId)))
            .CountDocumentsAsync();

        if (areFriends > 0)
        {
            return Results.BadRequest("You are already friends with this user");
        }

        // Check if a friend request is already pending (in either direction)
        var pendingFilter = Builders<FriendRequest>.Filter.Or(
            Builders<FriendRequest>.Filter.And(
                Builders<FriendRequest>.Filter.Eq("sender", senderObjId),
                Builders<FriendRequest>.Filter.Eq("recipient", recipientObjId)),
            Builders<FriendRequest>.Filter.And(
                Builders<FriendRequest>.Filter.Eq("sender", recipientObjId),
                Builders<FriendRequest>.Filter.Eq("recipient", senderObjId))
        );

        var existingRequest = await db.FriendRequests.Find(pendingFilter).FirstOrDefaultAsync();
        if (existingRequest != null)
        {
            return Results.BadRequest("A friend request already exist between you and this user");
        }

        var newRequest = new FriendRequest
        {
            Sender = userId,
            Recipient = id,
            Status = "pending"
        };

        await db.FriendRequests.InsertOneAsync(newRequest);

        return Results.Created("/api/users/friend-request", newRequest.Id);
    }

    private static async Task<IResult> AcceptFriendRequest(
        string id,
        HttpContext context,
        DatabaseService db)
    {
        if (context.Items["userId"] is not string userId)
        {
            return Results.Unauthorized();
        }

        if (string.IsNullOrEmpty(id))
        {
            return Results.BadRequest("No requestId provied");
        }

        if (!ObjectId.TryParse(id, out var requestId) || !ObjectId.TryParse(userId, out var userObjId))
        {
            return Results.BadRequest("Invalid id");
        }

        var existingRequest = await db.FriendRequests
            .Find(Builders<FriendRequest>.Filter.Eq("_id", requestId))
            .FirstOrDefaultAsync();

        if (existingRequest == null)
        {
            return Results.NotFound("No friend request founded");
        }

        if (existingRequest.Recipient != userId)
        {
            return Results.Unauthorized();
        }

        // Update the request status to accepted
        var update = Builders<FriendRequest>.Update.Set("status", "accepted");
        var updateResult = await db.FriendRequests.UpdateOneAsync(
            Builders<FriendRequest>.Filter.Eq("_id", requestId),
            update);

        if (updateResult.MatchedCount == 0)
        {
            return Results.StatusCode(StatusCodes.Status500InternalServerError);
        }

        // Add both users to each other's friends lists
        if (!ObjectId.TryParse(existingRequest.Sender, out var senderObjId))
        {
            return Results.StatusCode(StatusCodes.Status500InternalServerError);
        }

        await db.Users.UpdateOneAsync(
            Builders<User>.Filter.Eq("_id", userObjId),
            Builders<User>.Update.AddToSet("friends", existingRequest.Sender));

        await db.Users.UpdateOneAsync(
            Builders<User>.Filter.Eq("_id", senderObjId),
            Builders<User>.Update.AddToSet("friends", userId));

        return Results.Ok("Friend request accepted");
    }

    private static async Task<IResult> GetFriendRequests(HttpContext context, DatabaseService db)
    {
        if (context.Items["userId"] is not string userId)
        {
            return Results.Unauthorized();
        }

        if (!ObjectId.TryParse(userId, out var userObjId))
        {
            return Results.BadRequest("Invalid id");
        }

        // Pending requests (incoming) — user is the recipient
        var pendingPipeline = new[]
        {
            new BsonDocument("$match", new BsonDocument
            {
                { "recipient", userObjId },
                { "status", "pending" }
            }),
            new BsonDocument("$lookup", new BsonDocument
            {
                { "from", "users" },
                { "localField", "sender" },
                { "foreignField", "_id" },
                { "as", "user" }
            }),
            new BsonDocument("$unwind", "$user"),
            new BsonDocument("$project", new BsonDocument
            {
                { "_id", 0 },
                { "id", "$_id" },
                { "status", 1 },
                { "user", new BsonDocument
                    {
                        { "id", "$user._id" },
                        { "fullName", "$user.fullname" },
                        { "profilePic", "$user.profilepic" },
                        { "nativeLanguage", "$user.nativelanguage" },
                        { "learningLanguage", "$user.learninglanguage" }
                    }
                }
            })
        };

        // Accepted requests (user is the sender or recipient)
        var acceptedPipeline = new[]
        {
            new BsonDocument("$match", new BsonDocument
            {
                { "status", "accepted" },
                { "$or", new BsonArray
                    {
                        new BsonDocument("sender", userObjId)
                    }
                }
            }),
            new BsonDocument("$addFields", new BsonDocument
            {
                { "otherUserId", new BsonDocument("$cond", new BsonArray
                    {
                        new BsonDocument("$eq", new BsonArray { "$sender", userObjId }),
                        "$recipient",
                        "$sender"
                    })
                }
            }),
            new BsonDocument("$lookup", new BsonDocument
            {
                { "from", "users" },
                { "localField", "otherUserId" },
                { "foreignField", "_id" },
                { "as", "user" }
            }),
            new BsonDocument("$unwind", "$user"),
            new BsonDocument("$project", new BsonDocument
            {
                { "_id", 0 },
                { "id", "$_id" },
                { "status", 1 },
                { "user", new BsonDocument
                    {
                        { "id", "$user._id" },
                        { "fullName", "$user.fullname" },
                        { "profilePic", "$user.profilepic" }
                    }
                }
            })
        };

        var pendingResults = await db.FriendRequests.Aggregate<BsonDocument>(pendingPipeline).ToListAsync();
        var acceptedResults = await db.FriendRequests.Aggregate<BsonDocument>(acceptedPipeline).ToListAsync();

        return Results.Ok(new
        {
            PendingRequest = BsonJsonHelper.ToJsonList(pendingResults),
            AcceptedRequest = BsonJsonHelper.ToJsonList(acceptedResults)
        });
    }

    private static async Task<IResult> GetOutgoingFriendRequest(HttpContext context, DatabaseService db)
    {
        if (context.Items["userId"] is not string userId)
        {
            return Results.Unauthorized();
        }

        if (!ObjectId.TryParse(userId, out var userObjId))
        {
            return Results.BadRequest("Invalid id");
        }

        // Outgoing pending requests — user is the sender
        var pendingPipeline = new[]
        {
            new BsonDocument("$match", new BsonDocument
            {
                { "sender", userObjId },
                { "status", "pending" }
            }),
            new BsonDocument("$lookup", new BsonDocument
            {
                { "from", "users" },
                { "localField", "recipient" },
                { "foreignField", "_id" },
                { "as", "user" }
            }),
            new BsonDocument("$unwind", "$user"),
            new BsonDocument("$project", new BsonDocument
            {
                { "_id", 0 },
                { "id", "$_id" },
                { "status", 1 },
                { "user", new BsonDocument
                    {
                        { "id", "$user._id" },
                        { "fullName", "$user.fullname" },
                        { "profilePic", "$user.profilepic" },
                        { "nativeLanguage", "$user.nativelanguage" },
                        { "learningLanguage", "$user.learninglanguage" }
                    }
                }
            })
        };

        var pendingResults = await db.FriendRequests.Aggregate<BsonDocument>(pendingPipeline).ToListAsync();

        return Results.Ok(new
        {
            PendingRequest = BsonJsonHelper.ToJsonList(pendingResults)
        });
    }
}
