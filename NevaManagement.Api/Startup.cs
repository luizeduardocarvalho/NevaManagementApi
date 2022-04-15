using Microsoft.AspNetCore.Builder;
using Microsoft.AspNetCore.Hosting;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Configuration;
using Microsoft.Extensions.DependencyInjection;
using Microsoft.OpenApi.Models;
using NevaManagement.Domain.Interfaces.Repositories;
using NevaManagement.Domain.Interfaces.Services;
using NevaManagement.Infrastructure;
using NevaManagement.Infrastructure.Repositories;
using NevaManagement.Infrastructure.Services;
using Newtonsoft.Json.Serialization;
using System;

namespace NevaManagement.Api
{
    public class Startup
    {
        public Startup(IConfiguration configuration)
        {
            Configuration = configuration;
        }

        public IConfiguration Configuration { get; }

        public void ConfigureServices(IServiceCollection services)
        {
            services.AddCors();
            services.AddControllers();
            services.AddSwaggerGen(c =>
            {
                c.SwaggerDoc("v1", new OpenApiInfo { Title = "NevaManagement.Api", Version = "v1" });
            });

            services.AddDbContext<NevaManagementDbContext>(
                options =>
                    options.UseNpgsql(
                        GetConnectionString(),
                        x => x.MigrationsAssembly("NevaManagement.Infrastructure")));

            // Repositories
            services.AddTransient<IResearcherRepository, ResearcherRepository>();
            services.AddTransient<IProductRepository, ProductRepository>();
            services.AddTransient<ILocationRepository, LocationRepository>();
            services.AddTransient<IProductUsageService, ProductUsageService>();
            services.AddTransient<ILocationRepository, LocationRepository>();
            services.AddTransient<IOrganismRepository, OrganismRepository>();

            // Services
            services.AddTransient<IResearcherService, ResearcherService>();
            services.AddTransient<IProductService, ProductService>();
            services.AddTransient<IProductUsageRepository, ProductUsageRepository>();
            services.AddTransient<ILocationService, LocationService>();
            services.AddTransient<IOrganismService, OrganismService>();
        }

        public void Configure(IApplicationBuilder app, IWebHostEnvironment env)
        {
            app.UseDeveloperExceptionPage();
            app.UseSwagger();
            app.UseSwaggerUI(c => c.SwaggerEndpoint("/swagger/v1/swagger.json", "NevaManagement.Api v1"));

            app.UseHttpsRedirection();

            app.UseRouting();

            app.UseAuthorization();

            app.UseCors(
                options => options.AllowAnyOrigin().AllowAnyMethod().AllowAnyHeader()
            );

            app.UseEndpoints(endpoints =>
            {
                endpoints.MapControllers();
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
}
