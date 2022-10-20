using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace NevaManagement.Infrastructure.Migrations
{
    public partial class AddPatrimonyAndLocationToEquipment : Migration
    {
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.AddColumn<long>(
                name: "Location_Id",
                table: "Equipment",
                type: "bigint",
                nullable: true);

            migrationBuilder.AddColumn<string>(
                name: "Patrimony",
                table: "Equipment",
                type: "text",
                nullable: true);

            migrationBuilder.CreateIndex(
                name: "IX_Equipment_Location_Id",
                table: "Equipment",
                column: "Location_Id");

            migrationBuilder.AddForeignKey(
                name: "FK_Equipment_Location_Location_Id",
                table: "Equipment",
                column: "Location_Id",
                principalTable: "Location",
                principalColumn: "Id");
        }

        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropForeignKey(
                name: "FK_Equipment_Location_Location_Id",
                table: "Equipment");

            migrationBuilder.DropIndex(
                name: "IX_Equipment_Location_Id",
                table: "Equipment");

            migrationBuilder.DropColumn(
                name: "Location_Id",
                table: "Equipment");

            migrationBuilder.DropColumn(
                name: "Patrimony",
                table: "Equipment");
        }
    }
}
