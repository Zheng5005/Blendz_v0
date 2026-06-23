using Microsoft.AspNetCore.Http;

namespace ServerCSharp.Utils;

public static class CookieHelper
{
    public const string CookieName = "Blendz_Session";

    public static CookieOptions BuildSetCookieOptions()
    {
        return new CookieOptions
        {
            Expires = DateTimeOffset.UtcNow.AddDays(7),
            Path = "/",
            HttpOnly = true,
            Secure = true,
            SameSite = SameSiteMode.None
        };
    }

    public static CookieOptions BuildClearCookieOptions()
    {
        return new CookieOptions
        {
            Expires = DateTimeOffset.UtcNow.AddDays(-1),
            Path = "/",
            HttpOnly = true,
            Secure = true,
            SameSite = SameSiteMode.None
        };
    }

    public static void SetSessionCookie(HttpContext context, string token)
    {
        context.Response.Cookies.Append(CookieName, token, BuildSetCookieOptions());
    }

    public static void ClearSessionCookie(HttpContext context)
    {
        context.Response.Cookies.Append(CookieName, string.Empty, BuildClearCookieOptions());
    }
}
