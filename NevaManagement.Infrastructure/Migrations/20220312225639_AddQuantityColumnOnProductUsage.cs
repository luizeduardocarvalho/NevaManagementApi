using Microsoft.EntityFrameworkCore.Migrations;

namespace NevaManagement.Infrastructure.Migrations
{
    public partial class AddQuantityColumnOnProductUsage : Migration
    {
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.AddColumn<double>(
                name: "Quantity",
                table: "ProductUsage",
                type: "double precision",
                nullable: false,
                defaultValue: 0.0);
        }

        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropColumn(
                name: "Quantity",
                table: "ProductUsage");
        }
    }
}
