using Microsoft.EntityFrameworkCore.Migrations;

namespace NevaManagement.Infrastructure.Migrations
{
    public partial class CreateTables : Migration
    {
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropForeignKey(
                name: "FK_Location_Location_SubLocationId",
                table: "Location");

            migrationBuilder.DropForeignKey(
                name: "FK_Product_Location_LocationId",
                table: "Product");

            migrationBuilder.RenameColumn(
                name: "LocationId",
                table: "Product",
                newName: "Location_Id");

            migrationBuilder.RenameIndex(
                name: "IX_Product_LocationId",
                table: "Product",
                newName: "IX_Product_Location_Id");

            migrationBuilder.RenameColumn(
                name: "SubLocationId",
                table: "Location",
                newName: "Sub_Location_Id");

            migrationBuilder.RenameIndex(
                name: "IX_Location_SubLocationId",
                table: "Location",
                newName: "IX_Location_Sub_Location_Id");

            migrationBuilder.AddForeignKey(
                name: "FK_Location_Location_Sub_Location_Id",
                table: "Location",
                column: "Sub_Location_Id",
                principalTable: "Location",
                principalColumn: "Id",
                onDelete: ReferentialAction.Restrict);

            migrationBuilder.AddForeignKey(
                name: "FK_Product_Location_Location_Id",
                table: "Product",
                column: "Location_Id",
                principalTable: "Location",
                principalColumn: "Id",
                onDelete: ReferentialAction.Restrict);
        }

        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropForeignKey(
                name: "FK_Location_Location_Sub_Location_Id",
                table: "Location");

            migrationBuilder.DropForeignKey(
                name: "FK_Product_Location_Location_Id",
                table: "Product");

            migrationBuilder.RenameColumn(
                name: "Location_Id",
                table: "Product",
                newName: "LocationId");

            migrationBuilder.RenameIndex(
                name: "IX_Product_Location_Id",
                table: "Product",
                newName: "IX_Product_LocationId");

            migrationBuilder.RenameColumn(
                name: "Sub_Location_Id",
                table: "Location",
                newName: "SubLocationId");

            migrationBuilder.RenameIndex(
                name: "IX_Location_Sub_Location_Id",
                table: "Location",
                newName: "IX_Location_SubLocationId");

            migrationBuilder.AddForeignKey(
                name: "FK_Location_Location_SubLocationId",
                table: "Location",
                column: "SubLocationId",
                principalTable: "Location",
                principalColumn: "Id",
                onDelete: ReferentialAction.Restrict);

            migrationBuilder.AddForeignKey(
                name: "FK_Product_Location_LocationId",
                table: "Product",
                column: "LocationId",
                principalTable: "Location",
                principalColumn: "Id",
                onDelete: ReferentialAction.Restrict);
        }
    }
}
