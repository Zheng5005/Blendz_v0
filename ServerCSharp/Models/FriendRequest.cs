using MongoDB.Bson;
using MongoDB.Bson.Serialization.Attributes;

namespace ServerCSharp.Models;

public class FriendRequest
{
    [BsonId]
    [BsonRepresentation(BsonType.ObjectId)]
    public string? Id { get; set; }

    [BsonElement("sender")]
    [BsonRepresentation(BsonType.ObjectId)]
    public string Sender { get; set; } = string.Empty;

    [BsonElement("recipient")]
    [BsonRepresentation(BsonType.ObjectId)]
    public string Recipient { get; set; } = string.Empty;

    [BsonElement("status")]
    public string Status { get; set; } = string.Empty;
}
