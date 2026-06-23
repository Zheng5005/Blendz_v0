using MongoDB.Bson;
using MongoDB.Bson.Serialization.Attributes;

namespace ServerCSharp.Models;

public class UserCredentials
{
    [BsonId]
    [BsonRepresentation(BsonType.ObjectId)]
    public string? Id { get; set; }

    [BsonElement("email")]
    public string Email { get; set; } = string.Empty;

    [BsonElement("password")]
    public string Password { get; set; } = string.Empty;
}
