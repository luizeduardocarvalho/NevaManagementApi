using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace NevaManagement.Infrastructure.Migrations
{
    public partial class ChangePatromonyNameForPropertyNumber : Migration
    {
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.RenameColumn(
                name: "Patrimony",
                table: "Equipment",
                newName: "PropertyNumber");
        }

        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.RenameColumn(
                name: "PropertyNumber",
                table: "Equipment",
                newName: "Patrimony");
        }
    }
}
