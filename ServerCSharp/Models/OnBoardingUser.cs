using MongoDB.Bson.Serialization.Attributes;

namespace ServerCSharp.Models;

public class OnBoardingUser
{
    [BsonElement("fullname")]
    public string Fullname { get; set; } = string.Empty;

    [BsonElement("bio")]
    public string Bio { get; set; } = string.Empty;

    [BsonElement("nativelanguage")]
    public string NativeLanguage { get; set; } = string.Empty;

    [BsonElement("learninglanguage")]
    public string LearningLanguage { get; set; } = string.Empty;

    [BsonElement("location")]
    public string Location { get; set; } = string.Empty;

    [BsonElement("isonboarded")]
    public bool IsOnboarded { get; set; } = true;
}
