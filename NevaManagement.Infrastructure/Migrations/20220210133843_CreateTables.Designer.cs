﻿// <auto-generated />
using System;
using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Infrastructure;
using Microsoft.EntityFrameworkCore.Migrations;
using Microsoft.EntityFrameworkCore.Storage.ValueConversion;
using NevaManagement.Infrastructure;
using Npgsql.EntityFrameworkCore.PostgreSQL.Metadata;

namespace NevaManagement.Infrastructure.Migrations
{
    [DbContext(typeof(NevaManagementDbContext))]
    [Migration("20220210133843_CreateTables")]
    partial class CreateTables
    {
        protected override void BuildTargetModel(ModelBuilder modelBuilder)
        {
#pragma warning disable 612, 618
            modelBuilder
                .UseIdentityByDefaultColumns()
                .HasAnnotation("Relational:MaxIdentifierLength", 63)
                .HasAnnotation("ProductVersion", "5.0.0");

            modelBuilder.Entity("NevaManagement.Domain.Models.Location", b =>
                {
                    b.Property<long>("Id")
                        .ValueGeneratedOnAdd()
                        .HasColumnType("bigint")
                        .UseIdentityByDefaultColumn();

                    b.Property<string>("Description")
                        .HasColumnType("text");

                    b.Property<string>("Name")
                        .HasColumnType("text");

                    b.Property<long?>("Sub_Location_Id")
                        .HasColumnType("bigint");

                    b.HasKey("Id");

                    b.HasIndex("Sub_Location_Id");

                    b.ToTable("Location");
                });

            modelBuilder.Entity("NevaManagement.Domain.Models.Product", b =>
                {
                    b.Property<long>("Id")
                        .ValueGeneratedOnAdd()
                        .HasColumnType("bigint")
                        .UseIdentityByDefaultColumn();

                    b.Property<string>("Description")
                        .HasColumnType("text");

                    b.Property<long?>("Location_Id")
                        .HasColumnType("bigint");

                    b.Property<string>("Name")
                        .HasMaxLength(100)
                        .HasColumnType("character varying(100)");

                    b.Property<double>("Quantity")
                        .HasColumnType("double precision");

                    b.Property<string>("Unit")
                        .HasColumnType("text");

                    b.HasKey("Id");

                    b.HasIndex("Location_Id");

                    b.ToTable("Product");
                });

            modelBuilder.Entity("NevaManagement.Domain.Models.Location", b =>
                {
                    b.HasOne("NevaManagement.Domain.Models.Location", "SubLocation")
                        .WithMany()
                        .HasForeignKey("Sub_Location_Id");

                    b.Navigation("SubLocation");
                });

            modelBuilder.Entity("NevaManagement.Domain.Models.Product", b =>
                {
                    b.HasOne("NevaManagement.Domain.Models.Location", "Location")
                        .WithMany()
                        .HasForeignKey("Location_Id");

                    b.Navigation("Location");
                });
#pragma warning restore 612, 618
        }
    }
}
