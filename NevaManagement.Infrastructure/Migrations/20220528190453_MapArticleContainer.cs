using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace NevaManagement.Infrastructure.Migrations
{
    public partial class MapArticleContainer : Migration
    {
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropForeignKey(
                name: "FK_Article_Container_Container_Id",
                table: "Article");

            migrationBuilder.DropForeignKey(
                name: "FK_ArticleContainer_Article_Article_Id",
                table: "ArticleContainer");

            migrationBuilder.DropForeignKey(
                name: "FK_ArticleContainer_Container_Container_Id",
                table: "ArticleContainer");

            migrationBuilder.DropIndex(
                name: "IX_ArticleContainer_Article_Id",
                table: "ArticleContainer");

            migrationBuilder.DropIndex(
                name: "IX_Article_Container_Id",
                table: "Article");

            migrationBuilder.DropColumn(
                name: "Article_Id",
                table: "ArticleContainer");

            migrationBuilder.DropColumn(
                name: "Container_Id",
                table: "Article");

            migrationBuilder.RenameColumn(
                name: "Container_Id",
                table: "ArticleContainer",
                newName: "ArticleId");

            migrationBuilder.RenameIndex(
                name: "IX_ArticleContainer_Container_Id",
                table: "ArticleContainer",
                newName: "IX_ArticleContainer_ArticleId");

            migrationBuilder.AddColumn<long>(
                name: "ArticleID",
                table: "ArticleContainer",
                type: "bigint",
                nullable: false,
                defaultValue: 0L);

            migrationBuilder.AddColumn<long>(
                name: "ContainerId",
                table: "ArticleContainer",
                type: "bigint",
                nullable: false,
                defaultValue: 0L);

            migrationBuilder.CreateIndex(
                name: "IX_ArticleContainer_ContainerId",
                table: "ArticleContainer",
                column: "ContainerId");

            migrationBuilder.AddForeignKey(
                name: "FK_ArticleContainer_Article_ArticleId",
                table: "ArticleContainer",
                column: "ArticleId",
                principalTable: "Article",
                principalColumn: "Id");

            migrationBuilder.AddForeignKey(
                name: "FK_ArticleContainer_Container_ContainerId",
                table: "ArticleContainer",
                column: "ContainerId",
                principalTable: "Container",
                principalColumn: "Id",
                onDelete: ReferentialAction.Cascade);
        }

        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropForeignKey(
                name: "FK_ArticleContainer_Article_ArticleId",
                table: "ArticleContainer");

            migrationBuilder.DropForeignKey(
                name: "FK_ArticleContainer_Container_ContainerId",
                table: "ArticleContainer");

            migrationBuilder.DropIndex(
                name: "IX_ArticleContainer_ContainerId",
                table: "ArticleContainer");

            migrationBuilder.DropColumn(
                name: "ArticleID",
                table: "ArticleContainer");

            migrationBuilder.DropColumn(
                name: "ContainerId",
                table: "ArticleContainer");

            migrationBuilder.RenameColumn(
                name: "ArticleId",
                table: "ArticleContainer",
                newName: "Container_Id");

            migrationBuilder.RenameIndex(
                name: "IX_ArticleContainer_ArticleId",
                table: "ArticleContainer",
                newName: "IX_ArticleContainer_Container_Id");

            migrationBuilder.AddColumn<long>(
                name: "Article_Id",
                table: "ArticleContainer",
                type: "bigint",
                nullable: true);

            migrationBuilder.AddColumn<long>(
                name: "Container_Id",
                table: "Article",
                type: "bigint",
                nullable: true);

            migrationBuilder.CreateIndex(
                name: "IX_ArticleContainer_Article_Id",
                table: "ArticleContainer",
                column: "Article_Id");

            migrationBuilder.CreateIndex(
                name: "IX_Article_Container_Id",
                table: "Article",
                column: "Container_Id");

            migrationBuilder.AddForeignKey(
                name: "FK_Article_Container_Container_Id",
                table: "Article",
                column: "Container_Id",
                principalTable: "Container",
                principalColumn: "Id");

            migrationBuilder.AddForeignKey(
                name: "FK_ArticleContainer_Article_Article_Id",
                table: "ArticleContainer",
                column: "Article_Id",
                principalTable: "Article",
                principalColumn: "Id");

            migrationBuilder.AddForeignKey(
                name: "FK_ArticleContainer_Container_Container_Id",
                table: "ArticleContainer",
                column: "Container_Id",
                principalTable: "Container",
                principalColumn: "Id");
        }
    }
}
