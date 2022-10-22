using System;
using Microsoft.EntityFrameworkCore.Migrations;
using Npgsql.EntityFrameworkCore.PostgreSQL.Metadata;

#nullable disable

namespace NevaManagement.Infrastructure.Migrations
{
    public partial class InitialMigration : Migration
    {
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.CreateTable(
                name: "Article",
                columns: table => new
                {
                    Id = table.Column<long>(type: "bigint", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    Doi = table.Column<string>(type: "text", nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_Article", x => x.Id);
                });

            migrationBuilder.CreateTable(
                name: "Location",
                columns: table => new
                {
                    Id = table.Column<long>(type: "bigint", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    Name = table.Column<string>(type: "text", nullable: true),
                    Description = table.Column<string>(type: "text", nullable: true),
                    Sub_Location_Id = table.Column<long>(type: "bigint", nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_Location", x => x.Id);
                    table.ForeignKey(
                        name: "FK_Location_Location_Sub_Location_Id",
                        column: x => x.Sub_Location_Id,
                        principalTable: "Location",
                        principalColumn: "Id");
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
                        principalColumn: "Id");
                });

            migrationBuilder.CreateTable(
                name: "Researcher",
                columns: table => new
                {
                    Id = table.Column<long>(type: "bigint", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    Name = table.Column<string>(type: "character varying(80)", maxLength: 80, nullable: true),
                    Email = table.Column<string>(type: "text", nullable: true),
                    Password = table.Column<string>(type: "text", nullable: true),
                    Role = table.Column<string>(type: "text", nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_Researcher", x => x.Id);
                });

            migrationBuilder.CreateTable(
                name: "Equipment",
                columns: table => new
                {
                    Id = table.Column<long>(type: "bigint", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    Name = table.Column<string>(type: "text", nullable: true),
                    Description = table.Column<string>(type: "text", nullable: true),
                    Location_Id = table.Column<long>(type: "bigint", nullable: true),
                    PropertyNumber = table.Column<string>(type: "text", nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_Equipment", x => x.Id);
                    table.ForeignKey(
                        name: "FK_Equipment_Location_Location_Id",
                        column: x => x.Location_Id,
                        principalTable: "Location",
                        principalColumn: "Id");
                });

            migrationBuilder.CreateTable(
                name: "Product",
                columns: table => new
                {
                    Id = table.Column<long>(type: "bigint", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    Name = table.Column<string>(type: "character varying(100)", maxLength: 100, nullable: true),
                    Description = table.Column<string>(type: "text", nullable: true),
                    Location_Id = table.Column<long>(type: "bigint", nullable: true),
                    Quantity = table.Column<double>(type: "double precision", nullable: false),
                    Formula = table.Column<string>(type: "text", nullable: true),
                    Unit = table.Column<string>(type: "text", nullable: true),
                    ExpirationDate = table.Column<DateTimeOffset>(type: "timestamp with time zone", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_Product", x => x.Id);
                    table.ForeignKey(
                        name: "FK_Product_Location_Location_Id",
                        column: x => x.Location_Id,
                        principalTable: "Location",
                        principalColumn: "Id");
                });

            migrationBuilder.CreateTable(
                name: "Container",
                columns: table => new
                {
                    Id = table.Column<long>(type: "bigint", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    Creation_Date = table.Column<DateTimeOffset>(type: "timestamp with time zone", nullable: false),
                    Culture_Media = table.Column<string>(type: "text", nullable: true),
                    Description = table.Column<string>(type: "text", nullable: true),
                    Equipment_Id = table.Column<long>(type: "bigint", nullable: true),
                    Name = table.Column<string>(type: "character varying(80)", maxLength: 80, nullable: true),
                    Origin_Id = table.Column<long>(type: "bigint", nullable: true),
                    Researcher_Id = table.Column<long>(type: "bigint", nullable: true),
                    Sub_Container_Id = table.Column<long>(type: "bigint", nullable: true),
                    TransferDate = table.Column<DateTimeOffset>(type: "timestamp with time zone", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_Container", x => x.Id);
                    table.ForeignKey(
                        name: "FK_Container_Container_Sub_Container_Id",
                        column: x => x.Sub_Container_Id,
                        principalTable: "Container",
                        principalColumn: "Id");
                    table.ForeignKey(
                        name: "FK_Container_Equipment_Equipment_Id",
                        column: x => x.Equipment_Id,
                        principalTable: "Equipment",
                        principalColumn: "Id");
                    table.ForeignKey(
                        name: "FK_Container_Organism_Origin_Id",
                        column: x => x.Origin_Id,
                        principalTable: "Organism",
                        principalColumn: "Id");
                    table.ForeignKey(
                        name: "FK_Container_Researcher_Researcher_Id",
                        column: x => x.Researcher_Id,
                        principalTable: "Researcher",
                        principalColumn: "Id");
                });

            migrationBuilder.CreateTable(
                name: "EquipmentUsage",
                columns: table => new
                {
                    Id = table.Column<long>(type: "bigint", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    ResearcherId = table.Column<long>(type: "bigint", nullable: false),
                    EquipmentId = table.Column<long>(type: "bigint", nullable: false),
                    StartDate = table.Column<DateTimeOffset>(type: "timestamp with time zone", nullable: false),
                    EndDate = table.Column<DateTimeOffset>(type: "timestamp with time zone", nullable: false),
                    Description = table.Column<string>(type: "text", nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_EquipmentUsage", x => x.Id);
                    table.ForeignKey(
                        name: "FK_EquipmentUsage_Equipment_EquipmentId",
                        column: x => x.EquipmentId,
                        principalTable: "Equipment",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Cascade);
                    table.ForeignKey(
                        name: "FK_EquipmentUsage_Researcher_ResearcherId",
                        column: x => x.ResearcherId,
                        principalTable: "Researcher",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Cascade);
                });

            migrationBuilder.CreateTable(
                name: "ProductUsage",
                columns: table => new
                {
                    Id = table.Column<long>(type: "bigint", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    Researcher_Id = table.Column<long>(type: "bigint", nullable: true),
                    Product_Id = table.Column<long>(type: "bigint", nullable: true),
                    Usage_Date = table.Column<DateTimeOffset>(type: "timestamp with time zone", nullable: false),
                    Quantity = table.Column<double>(type: "double precision", nullable: false),
                    Description = table.Column<string>(type: "text", nullable: true)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_ProductUsage", x => x.Id);
                    table.ForeignKey(
                        name: "FK_ProductUsage_Product_Product_Id",
                        column: x => x.Product_Id,
                        principalTable: "Product",
                        principalColumn: "Id");
                    table.ForeignKey(
                        name: "FK_ProductUsage_Researcher_Researcher_Id",
                        column: x => x.Researcher_Id,
                        principalTable: "Researcher",
                        principalColumn: "Id");
                });

            migrationBuilder.CreateTable(
                name: "ArticleContainer",
                columns: table => new
                {
                    Id = table.Column<long>(type: "bigint", nullable: false)
                        .Annotation("Npgsql:ValueGenerationStrategy", NpgsqlValueGenerationStrategy.IdentityByDefaultColumn),
                    ArticleId = table.Column<long>(type: "bigint", nullable: false),
                    ContainerId = table.Column<long>(type: "bigint", nullable: false)
                },
                constraints: table =>
                {
                    table.PrimaryKey("PK_ArticleContainer", x => x.Id);
                    table.ForeignKey(
                        name: "FK_ArticleContainer_Article_ArticleId",
                        column: x => x.ArticleId,
                        principalTable: "Article",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Cascade);
                    table.ForeignKey(
                        name: "FK_ArticleContainer_Container_ContainerId",
                        column: x => x.ContainerId,
                        principalTable: "Container",
                        principalColumn: "Id",
                        onDelete: ReferentialAction.Cascade);
                });

            migrationBuilder.CreateIndex(
                name: "IX_ArticleContainer_ArticleId",
                table: "ArticleContainer",
                column: "ArticleId");

            migrationBuilder.CreateIndex(
                name: "IX_ArticleContainer_ContainerId",
                table: "ArticleContainer",
                column: "ContainerId");

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
                name: "IX_Equipment_Location_Id",
                table: "Equipment",
                column: "Location_Id");

            migrationBuilder.CreateIndex(
                name: "IX_EquipmentUsage_EquipmentId",
                table: "EquipmentUsage",
                column: "EquipmentId");

            migrationBuilder.CreateIndex(
                name: "IX_EquipmentUsage_ResearcherId",
                table: "EquipmentUsage",
                column: "ResearcherId");

            migrationBuilder.CreateIndex(
                name: "IX_Location_Sub_Location_Id",
                table: "Location",
                column: "Sub_Location_Id");

            migrationBuilder.CreateIndex(
                name: "IX_Organism_Origin_Id",
                table: "Organism",
                column: "Origin_Id");

            migrationBuilder.CreateIndex(
                name: "IX_Product_Location_Id",
                table: "Product",
                column: "Location_Id");

            migrationBuilder.CreateIndex(
                name: "IX_ProductUsage_Product_Id",
                table: "ProductUsage",
                column: "Product_Id");

            migrationBuilder.CreateIndex(
                name: "IX_ProductUsage_Researcher_Id",
                table: "ProductUsage",
                column: "Researcher_Id");
        }

        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropTable(
                name: "ArticleContainer");

            migrationBuilder.DropTable(
                name: "EquipmentUsage");

            migrationBuilder.DropTable(
                name: "ProductUsage");

            migrationBuilder.DropTable(
                name: "Article");

            migrationBuilder.DropTable(
                name: "Container");

            migrationBuilder.DropTable(
                name: "Product");

            migrationBuilder.DropTable(
                name: "Equipment");

            migrationBuilder.DropTable(
                name: "Organism");

            migrationBuilder.DropTable(
                name: "Researcher");

            migrationBuilder.DropTable(
                name: "Location");
        }
    }
}
