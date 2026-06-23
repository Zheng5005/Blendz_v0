using DotNetEnv;
using MongoDB.Driver;
using ServerCSharp.Models;

namespace ServerCSharp.Services;

public class DatabaseService
{
    public IMongoClient MongoClient { get; private set; } = null!;
    public IMongoCollection<User> Users { get; private set; } = null!;
    public IMongoCollection<FriendRequest> FriendRequests { get; private set; } = null!;

    public void Init()
    {
        Env.Load();

        var environment = Env.GetString("ENVIROMENT");
        var uri = environment == "development"
            ? Env.GetString("MONGODB_URI_DEV")
            : Env.GetString("MONGODB_URI");

        if (string.IsNullOrEmpty(uri))
        {
            throw new InvalidOperationException("MongoDB URI is not configured. Set MONGODB_URI or MONGODB_URI_DEV in .env");
        }

        var settings = MongoClientSettings.FromConnectionString(uri);
        MongoClient = new MongoClient(settings);

        // Force a connection so any URI errors are surfaced at startup
        MongoClient.GetDatabase("blendz").RunCommand<MongoDB.Bson.BsonDocument>(new MongoDB.Bson.BsonDocument("ping", 1));

        var database = MongoClient.GetDatabase("blendz");
        Users = database.GetCollection<User>("users");
        FriendRequests = database.GetCollection<FriendRequest>("friendRequests");

        Console.WriteLine("Connected to MongoDB");
    }
}
