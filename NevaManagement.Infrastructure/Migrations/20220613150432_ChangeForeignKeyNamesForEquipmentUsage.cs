using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace NevaManagement.Infrastructure.Migrations
{
    public partial class ChangeForeignKeyNamesForEquipmentUsage : Migration
    {
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropForeignKey(
                name: "FK_EquipmentUsage_Equipment_Equipment_Id",
                table: "EquipmentUsage");

            migrationBuilder.DropForeignKey(
                name: "FK_EquipmentUsage_Researcher_Researcher_Id",
                table: "EquipmentUsage");

            migrationBuilder.DropIndex(
                name: "IX_EquipmentUsage_Equipment_Id",
                table: "EquipmentUsage");

            migrationBuilder.DropIndex(
                name: "IX_EquipmentUsage_Researcher_Id",
                table: "EquipmentUsage");

            migrationBuilder.DropColumn(
                name: "Equipment_Id",
                table: "EquipmentUsage");

            migrationBuilder.DropColumn(
                name: "Researcher_Id",
                table: "EquipmentUsage");

            migrationBuilder.CreateIndex(
                name: "IX_EquipmentUsage_EquipmentId",
                table: "EquipmentUsage",
                column: "EquipmentId");

            migrationBuilder.CreateIndex(
                name: "IX_EquipmentUsage_ResearcherId",
                table: "EquipmentUsage",
                column: "ResearcherId");

            migrationBuilder.AddForeignKey(
                name: "FK_EquipmentUsage_Equipment_EquipmentId",
                table: "EquipmentUsage",
                column: "EquipmentId",
                principalTable: "Equipment",
                principalColumn: "Id",
                onDelete: ReferentialAction.Cascade);

            migrationBuilder.AddForeignKey(
                name: "FK_EquipmentUsage_Researcher_ResearcherId",
                table: "EquipmentUsage",
                column: "ResearcherId",
                principalTable: "Researcher",
                principalColumn: "Id",
                onDelete: ReferentialAction.Cascade);
        }

        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropForeignKey(
                name: "FK_EquipmentUsage_Equipment_EquipmentId",
                table: "EquipmentUsage");

            migrationBuilder.DropForeignKey(
                name: "FK_EquipmentUsage_Researcher_ResearcherId",
                table: "EquipmentUsage");

            migrationBuilder.DropIndex(
                name: "IX_EquipmentUsage_EquipmentId",
                table: "EquipmentUsage");

            migrationBuilder.DropIndex(
                name: "IX_EquipmentUsage_ResearcherId",
                table: "EquipmentUsage");

            migrationBuilder.AddColumn<long>(
                name: "Equipment_Id",
                table: "EquipmentUsage",
                type: "bigint",
                nullable: true);

            migrationBuilder.AddColumn<long>(
                name: "Researcher_Id",
                table: "EquipmentUsage",
                type: "bigint",
                nullable: true);

            migrationBuilder.CreateIndex(
                name: "IX_EquipmentUsage_Equipment_Id",
                table: "EquipmentUsage",
                column: "Equipment_Id");

            migrationBuilder.CreateIndex(
                name: "IX_EquipmentUsage_Researcher_Id",
                table: "EquipmentUsage",
                column: "Researcher_Id");

            migrationBuilder.AddForeignKey(
                name: "FK_EquipmentUsage_Equipment_Equipment_Id",
                table: "EquipmentUsage",
                column: "Equipment_Id",
                principalTable: "Equipment",
                principalColumn: "Id");

            migrationBuilder.AddForeignKey(
                name: "FK_EquipmentUsage_Researcher_Researcher_Id",
                table: "EquipmentUsage",
                column: "Researcher_Id",
                principalTable: "Researcher",
                principalColumn: "Id");
        }
    }
}
