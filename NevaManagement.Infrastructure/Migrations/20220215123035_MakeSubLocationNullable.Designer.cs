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
    [Migration("20220215123035_MakeSubLocationNullable")]
    partial class MakeSubLocationNullable
    {
        protected override void BuildTargetModel(ModelBuilder modelBuilder)
        {
#pragma warning disable 612, 618
            modelBuilder
                .UseIdentityByDefaultColumns()
                .HasAnnotation("Relational:MaxIdentifierLength", 63)
                .HasAnnotation("ProductVersion", "5.0.0");

            modelBuilder.Entity("NevaManagement.Domain.Models.Container", b =>
                {
                    b.Property<long>("Id")
                        .ValueGeneratedOnAdd()
                        .HasColumnType("bigint")
                        .UseIdentityByDefaultColumn();

                    b.Property<DateTimeOffset>("CreationDate")
                        .HasColumnType("timestamp with time zone")
                        .HasColumnName("Creation_Date");

                    b.Property<string>("CultureMedia")
                        .HasColumnType("text")
                        .HasColumnName("Culture_Media");

                    b.Property<string>("Description")
                        .HasColumnType("text");

                    b.Property<long?>("Equipment_Id")
                        .HasColumnType("bigint");

                    b.Property<string>("Name")
                        .HasMaxLength(80)
                        .HasColumnType("character varying(80)");

                    b.Property<long?>("Origin_Id")
                        .HasColumnType("bigint");

                    b.Property<long?>("Researcher_Id")
                        .HasColumnType("bigint");

                    b.Property<long?>("Sub_Container_Id")
                        .HasColumnType("bigint");

                    b.HasKey("Id");

                    b.HasIndex("Equipment_Id");

                    b.HasIndex("Origin_Id");

                    b.HasIndex("Researcher_Id");

                    b.HasIndex("Sub_Container_Id");

                    b.ToTable("Container");
                });

            modelBuilder.Entity("NevaManagement.Domain.Models.Equipment", b =>
                {
                    b.Property<long>("Id")
                        .ValueGeneratedOnAdd()
                        .HasColumnType("bigint")
                        .UseIdentityByDefaultColumn();

                    b.Property<string>("Description")
                        .HasColumnType("text");

                    b.Property<string>("Name")
                        .HasColumnType("text");

                    b.HasKey("Id");

                    b.ToTable("Equipment");
                });

            modelBuilder.Entity("NevaManagement.Domain.Models.EquipmentUsage", b =>
                {
                    b.Property<long>("Id")
                        .ValueGeneratedOnAdd()
                        .HasColumnType("bigint")
                        .UseIdentityByDefaultColumn();

                    b.Property<string>("Description")
                        .HasColumnType("text");

                    b.Property<long?>("Equipment_Id")
                        .HasColumnType("bigint");

                    b.Property<long?>("Researcher_Id")
                        .HasColumnType("bigint");

                    b.Property<DateTimeOffset>("UsageDate")
                        .HasColumnType("timestamp with time zone");

                    b.HasKey("Id");

                    b.HasIndex("Equipment_Id");

                    b.HasIndex("Researcher_Id");

                    b.ToTable("EquipmentUsage");
                });

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

            modelBuilder.Entity("NevaManagement.Domain.Models.Organism", b =>
                {
                    b.Property<long>("Id")
                        .ValueGeneratedOnAdd()
                        .HasColumnType("bigint")
                        .UseIdentityByDefaultColumn();

                    b.Property<DateTimeOffset>("CollectionDate")
                        .HasColumnType("timestamp with time zone")
                        .HasColumnName("Collection_Date");

                    b.Property<string>("CollectionLocation")
                        .HasColumnType("text")
                        .HasColumnName("Collection_Location");

                    b.Property<string>("Description")
                        .HasColumnType("text");

                    b.Property<DateTimeOffset>("IsolationDate")
                        .HasColumnType("timestamp with time zone")
                        .HasColumnName("Isolation_Date");

                    b.Property<string>("Name")
                        .HasMaxLength(100)
                        .HasColumnType("character varying(100)");

                    b.Property<string>("OriginPart")
                        .HasColumnType("text")
                        .HasColumnName("Origin_Part");

                    b.Property<long?>("Origin_Id")
                        .HasColumnType("bigint");

                    b.Property<string>("Type")
                        .HasColumnType("text");

                    b.HasKey("Id");

                    b.HasIndex("Origin_Id");

                    b.ToTable("Organism");
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

            modelBuilder.Entity("NevaManagement.Domain.Models.Researcher", b =>
                {
                    b.Property<long>("Id")
                        .ValueGeneratedOnAdd()
                        .HasColumnType("bigint")
                        .UseIdentityByDefaultColumn();

                    b.Property<string>("Email")
                        .HasMaxLength(100)
                        .HasColumnType("character varying(100)");

                    b.Property<string>("Name")
                        .HasMaxLength(80)
                        .HasColumnType("character varying(80)");

                    b.HasKey("Id");

                    b.ToTable("Researcher");
                });

            modelBuilder.Entity("NevaManagement.Domain.Models.Container", b =>
                {
                    b.HasOne("NevaManagement.Domain.Models.Equipment", "Equipment")
                        .WithMany()
                        .HasForeignKey("Equipment_Id");

                    b.HasOne("NevaManagement.Domain.Models.Organism", "Origin")
                        .WithMany()
                        .HasForeignKey("Origin_Id");

                    b.HasOne("NevaManagement.Domain.Models.Researcher", "Researcher")
                        .WithMany()
                        .HasForeignKey("Researcher_Id");

                    b.HasOne("NevaManagement.Domain.Models.Container", "SubContainer")
                        .WithMany()
                        .HasForeignKey("Sub_Container_Id");

                    b.Navigation("Equipment");

                    b.Navigation("Origin");

                    b.Navigation("Researcher");

                    b.Navigation("SubContainer");
                });

            modelBuilder.Entity("NevaManagement.Domain.Models.EquipmentUsage", b =>
                {
                    b.HasOne("NevaManagement.Domain.Models.Equipment", "Equipment")
                        .WithMany()
                        .HasForeignKey("Equipment_Id");

                    b.HasOne("NevaManagement.Domain.Models.Researcher", "Researcher")
                        .WithMany()
                        .HasForeignKey("Researcher_Id");

                    b.Navigation("Equipment");

                    b.Navigation("Researcher");
                });

            modelBuilder.Entity("NevaManagement.Domain.Models.Location", b =>
                {
                    b.HasOne("NevaManagement.Domain.Models.Location", "SubLocation")
                        .WithMany()
                        .HasForeignKey("Sub_Location_Id");

                    b.Navigation("SubLocation");
                });

            modelBuilder.Entity("NevaManagement.Domain.Models.Organism", b =>
                {
                    b.HasOne("NevaManagement.Domain.Models.Organism", "Origin")
                        .WithMany()
                        .HasForeignKey("Origin_Id");

                    b.Navigation("Origin");
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
