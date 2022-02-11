using Microsoft.EntityFrameworkCore.Migrations;

namespace NevaManagement.Infrastructure.Migrations
{
    public partial class FixCultureMediaColumnType : Migration
    {
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.RenameColumn(
                name: "Growth_Medium",
                table: "Container",
                newName: "Culture_Media");
        }

        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.RenameColumn(
                name: "Culture_Media",
                table: "Container",
                newName: "Growth_Medium");
        }
    }
}
