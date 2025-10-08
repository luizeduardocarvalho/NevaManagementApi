using System;
using Microsoft.EntityFrameworkCore.Migrations;
using Npgsql.EntityFrameworkCore.PostgreSQL.Metadata;

#nullable disable

namespace NevaManagement.Infrastructure.Migrations
{
    /// <inheritdoc />
    public partial class AddLaboratoryMultiTenancy : Migration
    {
        /// <inheritdoc />
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.AddColumn<long>(
                name: "LaboratoryId",
                table: "Researcher",
                type: "bigint",
                nullable: false,
                defaultValue: 0L);

            migrationBuilder.AddColumn<long>(
                name: "Laboratory_Id",
                table: "Researcher",
                type: "bigint",
                nullable: true);

            migrationBuilder.AddColumn<long>(
                name: "LaboratoryId",
                table: "ProductUsage",
                type: "bigint",
                nullable: false,
                defaultValue: 0L);

            migrationBuilder.AddColumn<long>(
                name: "Laboratory_Id",
                table: "ProductUsage",
                type: "bigint",
                nullable: true);

            migrationBuilder.AddColumn<long>(
                name: "LaboratoryId",
                table: "Product",
                type: "bigint",
                nullable: false,
                defaultValue: 0L);

            migrationBuilder.AddColumn<long>(
                name: "Laboratory_Id",
                table: "Product",
                type: "bigint",
                nullable: true);

            migrationBuilder.AddColumn<long>(
                name: "LaboratoryId",
                table: "Organism",
                type: "bigint",
                nullable: false,
                defaultValue: 0L);

            migrationBuilder.AddColumn<long>(
                name: "Laboratory_Id",
                table: "Organism",
                type: "bigint",
                nullable: true);

            migrationBuilder.AddColumn<long>(
                name: "LaboratoryId",
                table: "Location",
                type: "bigint",
                nullable: false,
                defaultValue: 0L);

            migrationBuilder.AddColumn<long>(
                name: "Laboratory_Id",
                table: "Location",
                type: "bigint",
                nullable: true);

            migrationBuilder.AddColumn<long>(
                name: "LaboratoryId",
                table: "EquipmentUsage",
                type: "bigint",
                nullable: false,
                defaultValue: 0L);

            migrationBuilder.AddColumn<long>(
                name: "Laboratory_Id",
                table: "EquipmentUsage",
                type: "bigint",
                nullable: true);

            migrationBuilder.AddColumn<long>(
                name: "LaboratoryId",
                table: "Equipment",
                type: "bigint",
                nullable: false,
                defaultValue: 0L);

            migrationBuilder.AddColumn<long>(
                name: "Laboratory_Id",
                table: "Equipment",
                type: "bigint",
                nullable: true);

            migrationBuilder.AddColumn<long>(
                name: "LaboratoryId",
                table: "Container",
                type: "bigint",
                nullable: false,
                defaultValue: 0L);

            migrationBuilder.AddColumn<long>(
                name: "Laboratory_Id",
                table: "Container",
                type: "bigint",
                nullable: true);

            migrationBuilder.CreateTable(
                name: "Laboratory",
                columns: table => new
                {
                    Id = table.Column<long>(type: "bigint", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    Name = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    Description = table.Column<string>(type: "text", nullable: true),
                    Address = table.Column<string>(type: "character varying(200)", maxLength: 200, nullable: true),
                    Phone = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: true),
                    Email = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    CreatedAt = table.Column<DateTimeOffset>(type: "timestamp with time zone", nullable: false),
                    IsActive = table.Column<bool>(type: "boolean", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_Laboratory", x => x.Id);
                });

            migrationBuilder.CreateTable(
                name: "LaboratoryInvitation",
                columns: table => new
                {
                    Id = table.Column<long>(type: "bigint", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    Laboratory_Id = table.Column<long>(type: "bigint", nullable: true),
                    LaboratoryId = table.Column<long>(type: "bigint", nullable: false),
                    InviteeEmail = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    InviteeName = table.Column<string>(type: "character varying(200)", maxLength: 200, nullable: true),
                    InvitationToken = table.Column<Guid>(type: "uuid", nullable: false),
                    CreatedAt = table.Column<DateTimeOffset>(type: "timestamp with time zone", nullable: false),
                    ExpiresAt = table.Column<DateTimeOffset>(type: "timestamp with time zone", nullable: false),
                    IsAccepted = table.Column<bool>(type: "boolean", nullable: false),
                    AcceptedAt = table.Column<DateTimeOffset>(type: "timestamp with time zone", nullable: true),
                    InvitedBy_Id = table.Column<long>(type: "bigint", nullable: true),
                    InvitedById = table.Column<long>(type: "bigint", nullable: false),
                    AcceptedBy_Id = table.Column<long>(type: "bigint", nullable: true),
                    AcceptedById = table.Column<long>(type: "bigint", nullable: true),
                    Role = table.Column<string>(type: "character varying(50)", maxLength: 50, nullable: true),
                    IsActive = table.Column<bool>(type: "boolean", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_LaboratoryInvitation", x => x.Id);
                    table.ForeignKey(
                        name: "FK_LaboratoryInvitation_Laboratory_Laboratory_Id",
                        column: x => x.Laboratory_Id,
                        principalTable: "Laboratory",
                        principalColumn: "Id");
                    table.ForeignKey(
                        name: "FK_LaboratoryInvitation_Researcher_AcceptedBy_Id",
                        column: x => x.AcceptedBy_Id,
                        principalTable: "Researcher",
                        principalColumn: "Id");
                    table.ForeignKey(
                        name: "FK_LaboratoryInvitation_Researcher_InvitedBy_Id",
                        column: x => x.InvitedBy_Id,
                        principalTable: "Researcher",
                        principalColumn: "Id");
                });

            migrationBuilder.CreateIndex(
                name: "IX_Researcher_Laboratory_Id",
                table: "Researcher",
                column: "Laboratory_Id");

            migrationBuilder.CreateIndex(
                name: "IX_ProductUsage_Laboratory_Id",
                table: "ProductUsage",
                column: "Laboratory_Id");

            migrationBuilder.CreateIndex(
                name: "IX_Product_Laboratory_Id",
                table: "Product",
                column: "Laboratory_Id");

            migrationBuilder.CreateIndex(
                name: "IX_Organism_Laboratory_Id",
                table: "Organism",
                column: "Laboratory_Id");

            migrationBuilder.CreateIndex(
                name: "IX_Location_Laboratory_Id",
                table: "Location",
                column: "Laboratory_Id");

            migrationBuilder.CreateIndex(
                name: "IX_EquipmentUsage_Laboratory_Id",
                table: "EquipmentUsage",
                column: "Laboratory_Id");

            migrationBuilder.CreateIndex(
                name: "IX_Equipment_Laboratory_Id",
                table: "Equipment",
                column: "Laboratory_Id");

            migrationBuilder.CreateIndex(
                name: "IX_Container_Laboratory_Id",
                table: "Container",
                column: "Laboratory_Id");

            migrationBuilder.CreateIndex(
                name: "IX_LaboratoryInvitation_AcceptedBy_Id",
                table: "LaboratoryInvitation",
                column: "AcceptedBy_Id");

            migrationBuilder.CreateIndex(
                name: "IX_LaboratoryInvitation_InvitedBy_Id",
                table: "LaboratoryInvitation",
                column: "InvitedBy_Id");

            migrationBuilder.CreateIndex(
                name: "IX_LaboratoryInvitation_Laboratory_Id",
                table: "LaboratoryInvitation",
                column: "Laboratory_Id");

            migrationBuilder.AddForeignKey(
                name: "FK_Container_Laboratory_Laboratory_Id",
                table: "Container",
                column: "Laboratory_Id",
                principalTable: "Laboratory",
                principalColumn: "Id");

            migrationBuilder.AddForeignKey(
                name: "FK_Equipment_Laboratory_Laboratory_Id",
                table: "Equipment",
                column: "Laboratory_Id",
                principalTable: "Laboratory",
                principalColumn: "Id");

            migrationBuilder.AddForeignKey(
                name: "FK_EquipmentUsage_Laboratory_Laboratory_Id",
                table: "EquipmentUsage",
                column: "Laboratory_Id",
                principalTable: "Laboratory",
                principalColumn: "Id");

            migrationBuilder.AddForeignKey(
                name: "FK_Location_Laboratory_Laboratory_Id",
                table: "Location",
                column: "Laboratory_Id",
                principalTable: "Laboratory",
                principalColumn: "Id");

            migrationBuilder.AddForeignKey(
                name: "FK_Organism_Laboratory_Laboratory_Id",
                table: "Organism",
                column: "Laboratory_Id",
                principalTable: "Laboratory",
                principalColumn: "Id");

            migrationBuilder.AddForeignKey(
                name: "FK_Product_Laboratory_Laboratory_Id",
                table: "Product",
                column: "Laboratory_Id",
                principalTable: "Laboratory",
                principalColumn: "Id");

            migrationBuilder.AddForeignKey(
                name: "FK_ProductUsage_Laboratory_Laboratory_Id",
                table: "ProductUsage",
                column: "Laboratory_Id",
                principalTable: "Laboratory",
                principalColumn: "Id");

            migrationBuilder.AddForeignKey(
                name: "FK_Researcher_Laboratory_Laboratory_Id",
                table: "Researcher",
                column: "Laboratory_Id",
                principalTable: "Laboratory",
                principalColumn: "Id");
        }

        /// <inheritdoc />
        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropForeignKey(
                name: "FK_Container_Laboratory_Laboratory_Id",
                table: "Container");

            migrationBuilder.DropForeignKey(
                name: "FK_Equipment_Laboratory_Laboratory_Id",
                table: "Equipment");

            migrationBuilder.DropForeignKey(
                name: "FK_EquipmentUsage_Laboratory_Laboratory_Id",
                table: "EquipmentUsage");

            migrationBuilder.DropForeignKey(
                name: "FK_Location_Laboratory_Laboratory_Id",
                table: "Location");

            migrationBuilder.DropForeignKey(
                name: "FK_Organism_Laboratory_Laboratory_Id",
                table: "Organism");

            migrationBuilder.DropForeignKey(
                name: "FK_Product_Laboratory_Laboratory_Id",
                table: "Product");

            migrationBuilder.DropForeignKey(
                name: "FK_ProductUsage_Laboratory_Laboratory_Id",
                table: "ProductUsage");

            migrationBuilder.DropForeignKey(
                name: "FK_Researcher_Laboratory_Laboratory_Id",
                table: "Researcher");

            migrationBuilder.DropTable(
                name: "LaboratoryInvitation");

            migrationBuilder.DropTable(
                name: "Laboratory");

            migrationBuilder.DropIndex(
                name: "IX_Researcher_Laboratory_Id",
                table: "Researcher");

            migrationBuilder.DropIndex(
                name: "IX_ProductUsage_Laboratory_Id",
                table: "ProductUsage");

            migrationBuilder.DropIndex(
                name: "IX_Product_Laboratory_Id",
                table: "Product");

            migrationBuilder.DropIndex(
                name: "IX_Organism_Laboratory_Id",
                table: "Organism");

            migrationBuilder.DropIndex(
                name: "IX_Location_Laboratory_Id",
                table: "Location");

            migrationBuilder.DropIndex(
                name: "IX_EquipmentUsage_Laboratory_Id",
                table: "EquipmentUsage");

            migrationBuilder.DropIndex(
                name: "IX_Equipment_Laboratory_Id",
                table: "Equipment");

            migrationBuilder.DropIndex(
                name: "IX_Container_Laboratory_Id",
                table: "Container");

            migrationBuilder.DropColumn(
                name: "LaboratoryId",
                table: "Researcher");

            migrationBuilder.DropColumn(
                name: "Laboratory_Id",
                table: "Researcher");

            migrationBuilder.DropColumn(
                name: "LaboratoryId",
                table: "ProductUsage");

            migrationBuilder.DropColumn(
                name: "Laboratory_Id",
                table: "ProductUsage");

            migrationBuilder.DropColumn(
                name: "LaboratoryId",
                table: "Product");

            migrationBuilder.DropColumn(
                name: "Laboratory_Id",
                table: "Product");

            migrationBuilder.DropColumn(
                name: "LaboratoryId",
                table: "Organism");

            migrationBuilder.DropColumn(
                name: "Laboratory_Id",
                table: "Organism");

            migrationBuilder.DropColumn(
                name: "LaboratoryId",
                table: "Location");

            migrationBuilder.DropColumn(
                name: "Laboratory_Id",
                table: "Location");

            migrationBuilder.DropColumn(
                name: "LaboratoryId",
                table: "EquipmentUsage");

            migrationBuilder.DropColumn(
                name: "Laboratory_Id",
                table: "EquipmentUsage");

            migrationBuilder.DropColumn(
                name: "LaboratoryId",
                table: "Equipment");

            migrationBuilder.DropColumn(
                name: "Laboratory_Id",
                table: "Equipment");

            migrationBuilder.DropColumn(
                name: "LaboratoryId",
                table: "Container");

            migrationBuilder.DropColumn(
                name: "Laboratory_Id",
                table: "Container");
        }
    }
}
