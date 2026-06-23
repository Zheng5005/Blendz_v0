using System.IdentityModel.Tokens.Jwt;
using System.Security.Claims;
using System.Security.Cryptography;
using System.Text;
using DotNetEnv;
using Microsoft.IdentityModel.Tokens;

namespace ServerCSharp.Services;

public class JwtService
{
    private readonly byte[] _keyBytes;

    public JwtService()
    {
        Env.Load();
        var secret = Env.GetString("JWT_SECRET");

        if (string.IsNullOrEmpty(secret))
        {
            throw new InvalidOperationException("JWT_SECRET is not configured. Set it in the .env file");
        }

        // HS256 requires a 256-bit (32-byte) key. Derive a key of exactly that length
        // from the configured secret so the server works regardless of secret length.
        _keyBytes = SHA256.HashData(Encoding.UTF8.GetBytes(secret));
    }

    public string GenerateJwt(string userId)
    {
        var key = new SymmetricSecurityKey(_keyBytes);
        var creds = new SigningCredentials(key, SecurityAlgorithms.HmacSha256);

        var token = new JwtSecurityToken(
            claims: new[]
            {
                new Claim("userId", userId)
            },
            expires: DateTime.UtcNow.AddHours(1),
            signingCredentials: creds
        );

        return new JwtSecurityTokenHandler().WriteToken(token);
    }

    public string? ValidateToken(string token)
    {
        try
        {
            var key = new SymmetricSecurityKey(_keyBytes);
            var tokenHandler = new JwtSecurityTokenHandler();

            tokenHandler.ValidateToken(token, new TokenValidationParameters
            {
                ValidateIssuer = false,
                ValidateAudience = false,
                ValidateLifetime = true,
                ValidateIssuerSigningKey = true,
                IssuerSigningKey = key,
                ClockSkew = TimeSpan.Zero
            }, out var validatedToken);

            var jwtToken = (JwtSecurityToken)validatedToken;
            var userId = jwtToken.Claims.FirstOrDefault(c => c.Type == "userId")?.Value;

            return userId;
        }
        catch
        {
            return null;
        }
    }
}
