using Microsoft.EntityFrameworkCore.Migrations;

namespace NevaManagement.Infrastructure.Migrations
{
    public partial class ChangeUsageDateColumnName : Migration
    {
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.RenameColumn(
                name: "UsageDate",
                table: "ProductUsage",
                newName: "Usage_Date");
        }

        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.RenameColumn(
                name: "Usage_Date",
                table: "ProductUsage",
                newName: "UsageDate");
        }
    }
}
