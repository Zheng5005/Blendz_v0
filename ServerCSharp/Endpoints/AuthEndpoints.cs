using MongoDB.Bson;
using MongoDB.Driver;
using ServerCSharp.Models;
using ServerCSharp.Services;
using ServerCSharp.Utils;

namespace ServerCSharp.Endpoints;

public static class AuthEndpoints
{
    public record SignupRequest(string FullName, string Email, string Password);
    public record LoginRequest(string Email, string Password);
    public record OnBoardingRequest(string FullName, string Bio, string NativeLanguage, string LearningLanguage, string Location);

    private static readonly System.Text.RegularExpressions.Regex EmailRegex = new(
        @"^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
        System.Text.RegularExpressions.RegexOptions.Compiled);

    public static void MapAuthEndpoints(this IEndpointRouteBuilder app)
    {
        app.MapPost("/api/auth/signup", Signup);
        app.MapPost("/api/auth/login", Login);
        app.MapPost("/api/auth/logout", Logout);
        app.MapPost("/api/auth/onboarding", OnBoard);
    }

    private static async Task<IResult> Signup(
        HttpContext context,
        SignupRequest input,
        DatabaseService db,
        StreamService streamService,
        JwtService jwtService)
    {
        if (input == null || string.IsNullOrEmpty(input.FullName) || string.IsNullOrEmpty(input.Email) || string.IsNullOrEmpty(input.Password))
        {
            return Results.BadRequest("Must have all required fields");
        }

        if (input.Password.Length < 6)
        {
            return Results.BadRequest("Password has to be at least 6 characters");
        }

        if (!EmailRegex.IsMatch(input.Email))
        {
            return Results.BadRequest("Invalid email");
        }

        var count = await db.Users.CountDocumentsAsync(Builders<User>.Filter.Eq("email", input.Email));
        if (count > 0)
        {
            return Results.BadRequest("Email already registered");
        }

        var hashedPassword = BCrypt.Net.BCrypt.HashPassword(input.Password);

        var newUser = new User
        {
            Fullname = input.FullName,
            Email = input.Email,
            Password = hashedPassword,
            ProfilePic = "https://avatar.iran.liara.run/public/12.png",
            Friends = new List<string>()
        };

        await db.Users.InsertOneAsync(newUser);
        var newUserId = newUser.Id!;

        try
        {
            await streamService.CreateStreamUserAsync(newUserId, newUser.Fullname, newUser.ProfilePic);
        }
        catch (Exception ex)
        {
            Console.WriteLine($"Stream user creation failed: {ex.Message}");
        }

        var token = jwtService.GenerateJwt(newUserId);
        context.Response.Cookies.Append(CookieHelper.CookieName, token, new CookieOptions
        {
            Expires = DateTimeOffset.UtcNow.AddDays(7),
            Path = "/",
            HttpOnly = true,
            Secure = true,
            SameSite = SameSiteMode.None
        });

        return Results.Created("/api/auth/signup", "Success");
    }

    private static async Task<IResult> Login(
        HttpContext context,
        LoginRequest input,
        DatabaseService db,
        JwtService jwtService)
    {
        if (input == null || string.IsNullOrEmpty(input.Email) || string.IsNullOrEmpty(input.Password))
        {
            return Results.BadRequest("All fields are required");
        }

        var filter = Builders<User>.Filter.Eq("email", input.Email);
        var projection = Builders<User>.Projection.Include("email").Include("password").Include("_id");
        var user = await db.Users.Find(filter).Project<User>(projection).FirstOrDefaultAsync();

        if (user == null || string.IsNullOrEmpty(user.Id))
        {
            return Results.Unauthorized();
        }

        if (string.IsNullOrEmpty(user.Password) || !BCrypt.Net.BCrypt.Verify(input.Password, user.Password))
        {
            return Results.Unauthorized();
        }

        var token = jwtService.GenerateJwt(user.Id);
        context.Response.Cookies.Append(CookieHelper.CookieName, token, new CookieOptions
        {
            Expires = DateTimeOffset.UtcNow.AddDays(7),
            Path = "/",
            HttpOnly = true,
            Secure = true,
            SameSite = SameSiteMode.None
        });

        return Results.Ok("Success");
    }

    private static IResult Logout(HttpContext context)
    {
        context.Response.Cookies.Append(CookieHelper.CookieName, string.Empty, new CookieOptions
        {
            Expires = DateTimeOffset.UtcNow.AddDays(-1),
            Path = "/",
            HttpOnly = true,
            Secure = true,
            SameSite = SameSiteMode.None
        });

        return Results.Ok("Logout");
    }

    private static async Task<IResult> OnBoard(
        HttpContext context,
        OnBoardingRequest input,
        DatabaseService db,
        StreamService streamService)
    {
        if (context.Items["userId"] is not string userId)
        {
            return Results.Unauthorized();
        }

        if (input == null ||
            string.IsNullOrEmpty(input.FullName) ||
            string.IsNullOrEmpty(input.Bio) ||
            string.IsNullOrEmpty(input.NativeLanguage) ||
            string.IsNullOrEmpty(input.LearningLanguage) ||
            string.IsNullOrEmpty(input.Location))
        {
            return Results.BadRequest("All fields are required");
        }

        if (!ObjectId.TryParse(userId, out var objectId))
        {
            return Results.BadRequest("Invalid id");
        }

        var update = Builders<User>.Update
            .Set("fullname", input.FullName)
            .Set("bio", input.Bio)
            .Set("nativelanguage", input.NativeLanguage)
            .Set("learninglanguage", input.LearningLanguage)
            .Set("location", input.Location)
            .Set("isonboarded", true);

        var result = await db.Users.UpdateOneAsync(Builders<User>.Filter.Eq("_id", objectId), update);

        if (result.MatchedCount == 0)
        {
            return Results.StatusCode(StatusCodes.Status500InternalServerError);
        }

        var user = await db.Users.Find(Builders<User>.Filter.Eq("_id", objectId)).FirstOrDefaultAsync();
        if (user == null || string.IsNullOrEmpty(user.Id))
        {
            return Results.StatusCode(StatusCodes.Status500InternalServerError);
        }

        try
        {
            await streamService.CreateStreamUserAsync(user.Id, user.Fullname, user.ProfilePic);
        }
        catch (Exception ex)
        {
            Console.WriteLine($"Stream user creation failed: {ex.Message}");
        }

        return Results.Ok("Success OnBoard");
    }
}
