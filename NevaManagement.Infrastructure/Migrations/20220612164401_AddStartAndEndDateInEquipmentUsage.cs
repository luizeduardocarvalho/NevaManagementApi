using System;
using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace NevaManagement.Infrastructure.Migrations
{
    public partial class AddStartAndEndDateInEquipmentUsage : Migration
    {
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.RenameColumn(
                name: "UsageDate",
                table: "EquipmentUsage",
                newName: "StartDate");

            migrationBuilder.AddColumn<DateTimeOffset>(
                name: "EndDate",
                table: "EquipmentUsage",
                type: "timestamp with time zone",
                nullable: false,
                defaultValue: new DateTimeOffset(new DateTime(1, 1, 1, 0, 0, 0, 0, DateTimeKind.Unspecified), new TimeSpan(0, 0, 0, 0, 0)));
        }

        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropColumn(
                name: "EndDate",
                table: "EquipmentUsage");

            migrationBuilder.RenameColumn(
                name: "StartDate",
                table: "EquipmentUsage",
                newName: "UsageDate");
        }
    }
}
