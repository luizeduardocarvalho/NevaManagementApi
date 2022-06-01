using System;
using Microsoft.EntityFrameworkCore.Migrations;
using Npgsql.EntityFrameworkCore.PostgreSQL.Metadata;

#nullable disable

namespace NevaManagement.Infrastructure.Migrations
{
    public partial class AddDoi : Migration
    {
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.AddColumn<DateTimeOffset>(
                name: "ExpirationDate",
                table: "Product",
                type: "timestamp with time zone",
                nullable: false,
                defaultValue: new DateTimeOffset(new DateTime(1, 1, 1, 0, 0, 0, 0, DateTimeKind.Unspecified), new TimeSpan(0, 0, 0, 0, 0)));

            migrationBuilder.CreateTable(
                name: "Article",
                columns: table => new
                {
                    Id = table.Column<long>(type: "bigint", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    Doi = table.Column<string>(type: "text", nullable: true),
                    Container_Id = table.Column<long>(type: "bigint", nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_Article", x => x.Id);
                    table.ForeignKey(
                        name: "FK_Article_Container_Container_Id",
                        column: x => x.Container_Id,
                        principalTable: "Container",
                        principalColumn: "Id");
                });

            migrationBuilder.CreateTable(
                name: "ArticleContainer",
                columns: table => new
                {
                    Id = table.Column<long>(type: "bigint", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    Article_Id = table.Column<long>(type: "bigint", nullable: true),
                    Container_Id = table.Column<long>(type: "bigint", nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_ArticleContainer", x => x.Id);
                    table.ForeignKey(
                        name: "FK_ArticleContainer_Article_Article_Id",
                        column: x => x.Article_Id,
                        principalTable: "Article",
                        principalColumn: "Id");
                    table.ForeignKey(
                        name: "FK_ArticleContainer_Container_Container_Id",
                        column: x => x.Container_Id,
                        principalTable: "Container",
                        principalColumn: "Id");
                });

            migrationBuilder.CreateIndex(
                name: "IX_Article_Container_Id",
                table: "Article",
                column: "Container_Id");

            migrationBuilder.CreateIndex(
                name: "IX_ArticleContainer_Article_Id",
                table: "ArticleContainer",
                column: "Article_Id");

            migrationBuilder.CreateIndex(
                name: "IX_ArticleContainer_Container_Id",
                table: "ArticleContainer",
                column: "Container_Id");
        }

        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "ArticleContainer");

            migrationBuilder.DropTable(
                name: "Article");

            migrationBuilder.DropColumn(
                name: "ExpirationDate",
                table: "Product");
        }
    }
}
