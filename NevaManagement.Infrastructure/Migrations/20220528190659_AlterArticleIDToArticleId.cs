using Microsoft.EntityFrameworkCore.Migrations;

#nullable disable

namespace NevaManagement.Infrastructure.Migrations
{
    public partial class AlterArticleIDToArticleId : Migration
    {
        protected override void Up(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropForeignKey(
                name: "FK_ArticleContainer_Article_ArticleId",
                table: "ArticleContainer");

            migrationBuilder.DropColumn(
                name: "ArticleID",
                table: "ArticleContainer");

            migrationBuilder.AlterColumn<long>(
                name: "ArticleId",
                table: "ArticleContainer",
                type: "bigint",
                nullable: false,
                defaultValue: 0L,
                oldClrType: typeof(long),
                oldType: "bigint",
                oldNullable: true);

            migrationBuilder.AddForeignKey(
                name: "FK_ArticleContainer_Article_ArticleId",
                table: "ArticleContainer",
                column: "ArticleId",
                principalTable: "Article",
                principalColumn: "Id",
                onDelete: ReferentialAction.Cascade);
        }

        protected override void Down(MigrationBuilder migrationBuilder)
        {
            migrationBuilder.DropForeignKey(
                name: "FK_ArticleContainer_Article_ArticleId",
                table: "ArticleContainer");

            migrationBuilder.AlterColumn<long>(
                name: "ArticleId",
                table: "ArticleContainer",
                type: "bigint",
                nullable: true,
                oldClrType: typeof(long),
                oldType: "bigint");

            migrationBuilder.AddColumn<long>(
                name: "ArticleID",
                table: "ArticleContainer",
                type: "bigint",
                nullable: false,
                defaultValue: 0L);

            migrationBuilder.AddForeignKey(
                name: "FK_ArticleContainer_Article_ArticleId",
                table: "ArticleContainer",
                column: "ArticleId",
                principalTable: "Article",
                principalColumn: "Id");
        }
    }
}
