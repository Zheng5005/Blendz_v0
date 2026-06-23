using MongoDB.Bson;
using MongoDB.Bson.Serialization.Attributes;

namespace ServerCSharp.Models;

[BsonIgnoreExtraElements]
public class User
{
    [BsonId]
    [BsonRepresentation(BsonType.ObjectId)]
    public string? Id { get; set; }

    [BsonElement("fullname")]
    [BsonIgnoreIfNull]
    public string Fullname { get; set; } = string.Empty;

    [BsonElement("email")]
    public string Email { get; set; } = string.Empty;

    [BsonElement("password")]
    public string Password { get; set; } = string.Empty;

    [BsonElement("bio")]
    [BsonIgnoreIfNull]
    public string Bio { get; set; } = string.Empty;

    [BsonElement("profilepic")]
    [BsonIgnoreIfNull]
    public string ProfilePic { get; set; } = string.Empty;

    [BsonElement("nativelanguage")]
    [BsonIgnoreIfNull]
    public string NativeLanguage { get; set; } = string.Empty;

    [BsonElement("learninglanguage")]
    [BsonIgnoreIfNull]
    public string LearningLanguage { get; set; } = string.Empty;

    [BsonElement("location")]
    [BsonIgnoreIfNull]
    public string Location { get; set; } = string.Empty;

    [BsonElement("isonboarded")]
    [BsonIgnoreIfDefault]
    public bool IsOnboarded { get; set; }

    [BsonElement("friends")]
    [BsonIgnoreIfNull]
    public List<string> Friends { get; set; } = new();
}
