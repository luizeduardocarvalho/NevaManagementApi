using Microsoft.AspNetCore.Cors.Infrastructure;
using NevaManagement.Api.Extensions;

namespace NevaManagement.Api;

public class Startup
{
    public Startup(IConfiguration configuration)
    {
        Configuration = configuration;
    }

    public IConfiguration Configuration { get; }

    public void ConfigureServices(IServiceCollection services)
    {
        services.AddMemoryCache();
        services.AddCors(opt => opt.AddPolicy("CorsPolicy", c =>
        {
            c.AllowAnyOrigin()
               .AllowAnyHeader()
               .AllowAnyMethod();
        }));
        services.AddControllers();
        services.AddSwaggerGen(c =>
        {
            c.SwaggerDoc("v1", new OpenApiInfo { Title = "NevaManagement.Api", Version = "v1" });
            c.AddSecurityDefinition("Bearer", new OpenApiSecurityScheme
            {
                Description =
                    "JWT Authorization header using the Bearer scheme. \r\n\r\n Enter 'Bearer' [space] and then your token in the text input below.\r\n\r\nExample: \"Bearer 12345abcdef\"",
                Name = "Authorization",
                In = ParameterLocation.Header,
                Type = SecuritySchemeType.ApiKey,
                Scheme = "Bearer"
            });

            c.AddSecurityRequirement(new OpenApiSecurityRequirement()
            {
                {
                    new OpenApiSecurityScheme
                    {
                        Reference = new OpenApiReference
                        {
                            Type = ReferenceType.SecurityScheme,
                            Id = "Bearer"
                        },
                        Scheme = "oauth2",
                        Name = "Bearer",
                        In = ParameterLocation.Header,

                    },
                    new List<string>()
                }
            });
        });

        services.AddDbContext<NevaManagementDbContext>(
            options =>
                options.UseNpgsql(
                    GetConnectionString(),
                    x => x.MigrationsAssembly("NevaManagement.Infrastructure")));

        services.Configure<Settings>(Configuration.GetSection("Settings"));


        var settings = Environment.GetEnvironmentVariable("Settings")
                       ?? Configuration.GetSection("Settings").GetSection("Secret").Value;

        services
            .AddAuthentication(x =>
            {
                x.DefaultAuthenticateScheme = JwtBearerDefaults.AuthenticationScheme;
                x.DefaultChallengeScheme = JwtBearerDefaults.AuthenticationScheme;
            })
            .AddJwtBearer(x =>
            {
                x.RequireHttpsMetadata = false;
                x.SaveToken = true;
                x.TokenValidationParameters = new TokenValidationParameters
                {
                    ValidateIssuerSigningKey = true,
                    IssuerSigningKey = new SymmetricSecurityKey(
                        Encoding.ASCII.GetBytes(settings)),
                    ValidateIssuer = false,
                    ValidateAudience = false
                };
            });

        services.ConfigureBuilders();
        services.ConfigureRepositories();
        services.ConfigureServices();

        services.ConfigureAutoMapper();
    }

    public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
    {
        app.UseDeveloperExceptionPage();
        app.UseSwagger();
        app.UseSwaggerUI(c => c.SwaggerEndpoint("/swagger/v1/swagger.json", "NevaManagement.Api v1"));

        app.UseHttpsRedirection();

        app.UseRouting();

        app.UseAuthentication();
        app.UseAuthorization();

        app.UseCors("CorsPolicy");

        app.InitializeCache();

        app.UseMiddleware<ErrorHandlerMiddleware>();

        app.UseEndpoints(endpoints =>
        {
            endpoints.MapControllers().RequireAuthorization();
        });
    }

    public string GetConnectionString()
    {
        var uriString = Environment.GetEnvironmentVariable("DATABASE_URL")
                            ?? Configuration.GetConnectionString("DATABASE_URL");
        var uri = new Uri(uriString);
        var db = uri.AbsolutePath.Trim('/');
        var user = uri.UserInfo.Split(':')[0];
        var passwd = uri.UserInfo.Split(':')[1];
        var port = uri.Port > 0 ? uri.Port : 5432;
        var connStr = string.Format("Server={0};Database={1};User Id={2};Password={3};Port={4};SSL Mode=Require;Trust Server Certificate=True;",
            uri.Host, db, user, passwd, port);
        return connStr;
    }
}
