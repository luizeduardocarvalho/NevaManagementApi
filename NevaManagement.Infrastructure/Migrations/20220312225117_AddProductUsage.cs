using System;
using Microsoft.EntityFrameworkCore.Migrations;
using Npgsql.EntityFrameworkCore.PostgreSQL.Metadata;

namespace NevaManagement.Infrastructure.Migrations
{
    public partial class AddProductUsage : Migration
    {
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.CreateTable(
                name: "ProductUsage",
                columns: table => new
                {
                    Id = table.Column<long>(type: "bigint", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    Researcher_Id = table.Column<long>(type: "bigint", nullable: true),
                    Product_Id = table.Column<long>(type: "bigint", nullable: true),
                    UsageDate = table.Column<DateTimeOffset>(type: "timestamp with time zone", nullable: false),
                    Description = table.Column<string>(type: "text", nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_ProductUsage", x => x.Id);
                    table.ForeignKey(
                        name: "FK_ProductUsage_Product_Product_Id",
                        column: x => x.Product_Id,
                        principalTable: "Product",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Restrict);
                    table.ForeignKey(
                        name: "FK_ProductUsage_Researcher_Researcher_Id",
                        column: x => x.Researcher_Id,
                        principalTable: "Researcher",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Restrict);
                });

            migrationBuilder.CreateIndex(
                name: "IX_ProductUsage_Product_Id",
                table: "ProductUsage",
                column: "Product_Id");

            migrationBuilder.CreateIndex(
                name: "IX_ProductUsage_Researcher_Id",
                table: "ProductUsage",
                column: "Researcher_Id");
        }

        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "ProductUsage");
        }
    }
}
