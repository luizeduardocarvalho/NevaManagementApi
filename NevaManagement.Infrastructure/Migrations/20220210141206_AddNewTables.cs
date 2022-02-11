using System;
using Microsoft.EntityFrameworkCore.Migrations;
using Npgsql.EntityFrameworkCore.PostgreSQL.Metadata;

namespace NevaManagement.Infrastructure.Migrations
{
    public partial class AddNewTables : Migration
    {
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.CreateTable(
                name: "Equipment",
                columns: table => new
                {
                    Id = table.Column<long>(type: "bigint", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    Name = table.Column<string>(type: "text", nullable: true),
                    Description = table.Column<string>(type: "text", nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_Equipment", x => x.Id);
                });

            migrationBuilder.CreateTable(
                name: "Organism",
                columns: table => new
                {
                    Id = table.Column<long>(type: "bigint", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    Name = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    Type = table.Column<string>(type: "text", nullable: true),
                    Description = table.Column<string>(type: "text", nullable: true),
                    Collection_Date = table.Column<DateTimeOffset>(type: "timestamp with time zone", nullable: false),
                    Collection_Location = table.Column<string>(type: "text", nullable: true),
                    Isolation_Date = table.Column<DateTimeOffset>(type: "timestamp with time zone", nullable: false),
                    Origin_Id = table.Column<long>(type: "bigint", nullable: true),
                    Origin_Part = table.Column<string>(type: "text", nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_Organism", x => x.Id);
                    table.ForeignKey(
                        name: "FK_Organism_Organism_Origin_Id",
                        column: x => x.Origin_Id,
                        principalTable: "Organism",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Restrict);
                });

            migrationBuilder.CreateTable(
                name: "Researcher",
                columns: table => new
                {
                    Id = table.Column<long>(type: "bigint", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    Name = table.Column<string>(type: "character varying(80)", maxLength: 80, nullable: true),
                    Email = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_Researcher", x => x.Id);
                });

            migrationBuilder.CreateTable(
                name: "Container",
                columns: table => new
                {
                    Id = table.Column<long>(type: "bigint", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    Name = table.Column<string>(type: "character varying(80)", maxLength: 80, nullable: true),
                    Description = table.Column<string>(type: "text", nullable: true),
                    Growth_Medium = table.Column<string>(type: "text", nullable: true),
                    Equipment_Id = table.Column<long>(type: "bigint", nullable: true),
                    Researcher_Id = table.Column<long>(type: "bigint", nullable: true),
                    Origin_Id = table.Column<long>(type: "bigint", nullable: true),
                    Sub_Container_Id = table.Column<long>(type: "bigint", nullable: true),
                    Creation_Date = table.Column<DateTimeOffset>(type: "timestamp with time zone", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_Container", x => x.Id);
                    table.ForeignKey(
                        name: "FK_Container_Container_Sub_Container_Id",
                        column: x => x.Sub_Container_Id,
                        principalTable: "Container",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Restrict);
                    table.ForeignKey(
                        name: "FK_Container_Equipment_Equipment_Id",
                        column: x => x.Equipment_Id,
                        principalTable: "Equipment",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Restrict);
                    table.ForeignKey(
                        name: "FK_Container_Organism_Origin_Id",
                        column: x => x.Origin_Id,
                        principalTable: "Organism",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Restrict);
                    table.ForeignKey(
                        name: "FK_Container_Researcher_Researcher_Id",
                        column: x => x.Researcher_Id,
                        principalTable: "Researcher",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Restrict);
                });

            migrationBuilder.CreateTable(
                name: "EquipmentUsage",
                columns: table => new
                {
                    Id = table.Column<long>(type: "bigint", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    Researcher_Id = table.Column<long>(type: "bigint", nullable: true),
                    Equipment_Id = table.Column<long>(type: "bigint", nullable: true),
                    UsageDate = table.Column<DateTimeOffset>(type: "timestamp with time zone", nullable: false),
                    Description = table.Column<string>(type: "text", nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_EquipmentUsage", x => x.Id);
                    table.ForeignKey(
                        name: "FK_EquipmentUsage_Equipment_Equipment_Id",
                        column: x => x.Equipment_Id,
                        principalTable: "Equipment",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Restrict);
                    table.ForeignKey(
                        name: "FK_EquipmentUsage_Researcher_Researcher_Id",
                        column: x => x.Researcher_Id,
                        principalTable: "Researcher",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Restrict);
                });

            migrationBuilder.CreateIndex(
                name: "IX_Container_Equipment_Id",
                table: "Container",
                column: "Equipment_Id");

            migrationBuilder.CreateIndex(
                name: "IX_Container_Origin_Id",
                table: "Container",
                column: "Origin_Id");

            migrationBuilder.CreateIndex(
                name: "IX_Container_Researcher_Id",
                table: "Container",
                column: "Researcher_Id");

            migrationBuilder.CreateIndex(
                name: "IX_Container_Sub_Container_Id",
                table: "Container",
                column: "Sub_Container_Id");

            migrationBuilder.CreateIndex(
                name: "IX_EquipmentUsage_Equipment_Id",
                table: "EquipmentUsage",
                column: "Equipment_Id");

            migrationBuilder.CreateIndex(
                name: "IX_EquipmentUsage_Researcher_Id",
                table: "EquipmentUsage",
                column: "Researcher_Id");

            migrationBuilder.CreateIndex(
                name: "IX_Organism_Origin_Id",
                table: "Organism",
                column: "Origin_Id");
        }

        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "Container");

            migrationBuilder.DropTable(
                name: "EquipmentUsage");

            migrationBuilder.DropTable(
                name: "Organism");

            migrationBuilder.DropTable(
                name: "Equipment");

            migrationBuilder.DropTable(
                name: "Researcher");
        }
    }
}
