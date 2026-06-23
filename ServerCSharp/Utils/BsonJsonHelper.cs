using System.Text.Json;
using MongoDB.Bson;

namespace ServerCSharp.Utils;

/// <summary>
/// Converts BsonDocument results from MongoDB aggregation pipelines into
/// a JSON-serializable structure. Each BsonValue is mapped to its closest
/// .NET equivalent (ObjectId → string, etc.) so the response matches what
/// the original Go server emits via json.Encoder.
/// </summary>
public static class BsonJsonHelper
{
    public static List<Dictionary<string, object?>> ToJsonList(IEnumerable<BsonDocument> documents)
    {
        return documents.Select(ToJsonMap).ToList();
    }

    public static Dictionary<string, object?> ToJsonMap(BsonDocument document)
    {
        var result = new Dictionary<string, object?>();
        foreach (var element in document)
        {
            result[element.Name] = ConvertValue(element.Value);
        }
        return result;
    }

    private static object? ConvertValue(BsonValue value)
    {
        return value.BsonType switch
        {
            BsonType.Document => ToJsonMap(value.AsBsonDocument),
            BsonType.Array => value.AsBsonArray.Select(ConvertValue).ToList(),
            BsonType.ObjectId => value.AsObjectId.ToString(),
            BsonType.String => value.AsString,
            BsonType.Int32 => value.AsInt32,
            BsonType.Int64 => value.AsInt64,
            BsonType.Double => value.AsDouble,
            BsonType.Boolean => value.AsBoolean,
            BsonType.Null => null,
            BsonType.DateTime => value.ToUniversalTime(),
            _ => value.ToString()
        };
    }
}
