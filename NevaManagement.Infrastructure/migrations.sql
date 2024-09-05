﻿CREATE TABLE IF NOT EXISTS "__EFMigrationsHistory" (
    "MigrationId" character varying(150) NOT NULL,
    "ProductVersion" character varying(32) NOT NULL,
    CONSTRAINT "PK___EFMigrationsHistory" PRIMARY KEY ("MigrationId")
);

DO $EF$
BEGIN
    IF NOT EXISTS(SELECT 1 FROM "__EFMigrationsHistory" WHERE "MigrationId" = '20221022152249_InitialMigration') THEN
        -- Article table
        CREATE TABLE IF NOT EXISTS "Article" (
            "Id" bigint GENERATED BY DEFAULT AS IDENTITY,
            "Doi" text NULL,
            CONSTRAINT "PK_Article" PRIMARY KEY ("Id")
        );

        -- Location table
        CREATE TABLE IF NOT EXISTS "Location" (
            "Id" bigint GENERATED BY DEFAULT AS IDENTITY,
            "Name" text NULL,
            "Description" text NULL,
            "Sub_Location_Id" bigint NULL,
            CONSTRAINT "PK_Location" PRIMARY KEY ("Id"),
            CONSTRAINT "FK_Location_Location_Sub_Location_Id" FOREIGN KEY ("Sub_Location_Id") REFERENCES "Location" ("Id")
        );

        -- Organism table
        CREATE TABLE IF NOT EXISTS "Organism" (
            "Id" bigint GENERATED BY DEFAULT AS IDENTITY,
            "Name" character varying(100) NULL,
            "Type" text NULL,
            "Description" text NULL,
            "Collection_Date" timestamp with time zone NOT NULL,
            "Collection_Location" text NULL,
            "Isolation_Date" timestamp with time zone NOT NULL,
            "Origin_Id" bigint NULL,
            "Origin_Part" text NULL,
            CONSTRAINT "PK_Organism" PRIMARY KEY ("Id"),
            CONSTRAINT "FK_Organism_Organism_Origin_Id" FOREIGN KEY ("Origin_Id") REFERENCES "Organism" ("Id")
        );

        -- Researcher table
        CREATE TABLE IF NOT EXISTS "Researcher" (
            "Id" bigint GENERATED BY DEFAULT AS IDENTITY,
            "Name" character varying(80) NULL,
            "Email" text NULL,
            "Password" text NULL,
            "Role" text NULL,
            CONSTRAINT "PK_Researcher" PRIMARY KEY ("Id")
        );

        -- Equipment table
        CREATE TABLE IF NOT EXISTS "Equipment" (
            "Id" bigint GENERATED BY DEFAULT AS IDENTITY,
            "Name" text NULL,
            "Description" text NULL,
            "Location_Id" bigint NULL,
            "PropertyNumber" text NULL,
            CONSTRAINT "PK_Equipment" PRIMARY KEY ("Id"),
            CONSTRAINT "FK_Equipment_Location_Location_Id" FOREIGN KEY ("Location_Id") REFERENCES "Location" ("Id")
        );

        -- Product table
        CREATE TABLE IF NOT EXISTS "Product" (
            "Id" bigint GENERATED BY DEFAULT AS IDENTITY,
            "Name" character varying(100) NULL,
            "Description" text NULL,
            "Location_Id" bigint NULL,
            "Quantity" double precision NOT NULL,
            "Formula" text NULL,
            "Unit" text NULL,
            "ExpirationDate" timestamp with time zone NOT NULL,
            CONSTRAINT "PK_Product" PRIMARY KEY ("Id"),
            CONSTRAINT "FK_Product_Location_Location_Id" FOREIGN KEY ("Location_Id") REFERENCES "Location" ("Id")
        );

        -- Container table
        CREATE TABLE IF NOT EXISTS "Container" (
            "Id" bigint GENERATED BY DEFAULT AS IDENTITY,
            "Creation_Date" timestamp with time zone NOT NULL,
            "Culture_Media" text NULL,
            "Description" text NULL,
            "Equipment_Id" bigint NULL,
            "Name" character varying(80) NULL,
            "Origin_Id" bigint NULL,
            "Researcher_Id" bigint NULL,
            "Sub_Container_Id" bigint NULL,
            "TransferDate" timestamp with time zone NOT NULL,
            CONSTRAINT "PK_Container" PRIMARY KEY ("Id"),
            CONSTRAINT "FK_Container_Container_Sub_Container_Id" FOREIGN KEY ("Sub_Container_Id") REFERENCES "Container" ("Id"),
            CONSTRAINT "FK_Container_Equipment_Equipment_Id" FOREIGN KEY ("Equipment_Id") REFERENCES "Equipment" ("Id"),
            CONSTRAINT "FK_Container_Organism_Origin_Id" FOREIGN KEY ("Origin_Id") REFERENCES "Organism" ("Id"),
            CONSTRAINT "FK_Container_Researcher_Researcher_Id" FOREIGN KEY ("Researcher_Id") REFERENCES "Researcher" ("Id")
        );

        -- ArticleContainer table
        CREATE TABLE IF NOT EXISTS "ArticleContainer" (
            "Id" bigint GENERATED BY DEFAULT AS IDENTITY,
            "ArticleId" bigint NOT NULL,
            "ContainerId" bigint NOT NULL,
            CONSTRAINT "PK_ArticleContainer" PRIMARY KEY ("Id"),
            CONSTRAINT "FK_ArticleContainer_Article_ArticleId" FOREIGN KEY ("ArticleId") REFERENCES "Article" ("Id") ON DELETE CASCADE,
            CONSTRAINT "FK_ArticleContainer_Container_ContainerId" FOREIGN KEY ("ContainerId") REFERENCES "Container" ("Id") ON DELETE CASCADE
        );

        -- EquipmentUsage table
        CREATE TABLE IF NOT EXISTS "EquipmentUsage" (
            "Id" bigint GENERATED BY DEFAULT AS IDENTITY,
            "ResearcherId" bigint NOT NULL,
            "EquipmentId" bigint NOT NULL,
            "StartDate" timestamp with time zone NOT NULL,
            "EndDate" timestamp with time zone NOT NULL,
            "Description" text NULL,
            CONSTRAINT "PK_EquipmentUsage" PRIMARY KEY ("Id"),
            CONSTRAINT "FK_EquipmentUsage_Equipment_EquipmentId" FOREIGN KEY ("EquipmentId") REFERENCES "Equipment" ("Id") ON DELETE CASCADE,
            CONSTRAINT "FK_EquipmentUsage_Researcher_ResearcherId" FOREIGN KEY ("ResearcherId") REFERENCES "Researcher" ("Id") ON DELETE CASCADE
        );

        -- ProductUsage table
        CREATE TABLE IF NOT EXISTS "ProductUsage" (
            "Id" bigint GENERATED BY DEFAULT AS IDENTITY,
            "Researcher_Id" bigint NULL,
            "Product_Id" bigint NULL,
            "Usage_Date" timestamp with time zone NOT NULL,
            "Quantity" double precision NOT NULL,
            "Description" text NULL,
            CONSTRAINT "PK_ProductUsage" PRIMARY KEY ("Id"),
            CONSTRAINT "FK_ProductUsage_Product_Product_Id" FOREIGN KEY ("Product_Id") REFERENCES "Product" ("Id"),
            CONSTRAINT "FK_ProductUsage_Researcher_Researcher_Id" FOREIGN KEY ("Researcher_Id") REFERENCES "Researcher" ("Id")
        );

        -- Index creation after all relevant tables have been created
        CREATE INDEX IF NOT EXISTS "IX_ArticleContainer_ArticleId" ON "ArticleContainer" ("ArticleId");
        CREATE INDEX IF NOT EXISTS "IX_ArticleContainer_ContainerId" ON "ArticleContainer" ("ContainerId");
        CREATE INDEX IF NOT EXISTS "IX_Container_Equipment_Id" ON "Container" ("Equipment_Id");
        CREATE INDEX IF NOT EXISTS "IX_Container_Origin_Id" ON "Container" ("Origin_Id");
        CREATE INDEX IF NOT EXISTS "IX_Container_Researcher_Id" ON "Container" ("Researcher_Id");
        CREATE INDEX IF NOT EXISTS "IX_Container_Sub_Container_Id" ON "Container" ("Sub_Container_Id");
        CREATE INDEX IF NOT EXISTS "IX_Equipment_Location_Id" ON "Equipment" ("Location_Id");
        CREATE INDEX IF NOT EXISTS "IX_EquipmentUsage_EquipmentId" ON "EquipmentUsage" ("EquipmentId");
        CREATE INDEX IF NOT EXISTS "IX_EquipmentUsage_ResearcherId" ON "EquipmentUsage" ("ResearcherId");
        CREATE INDEX IF NOT EXISTS "IX_Location_Sub_Location_Id" ON "Location" ("Sub_Location_Id");
        CREATE INDEX IF NOT EXISTS "IX_Organism_Origin_Id" ON "Organism" ("Origin_Id");
        CREATE INDEX IF NOT EXISTS "IX_Product_Location_Id" ON "Product" ("Location_Id");
        CREATE INDEX IF NOT EXISTS "IX_ProductUsage_Product_Id" ON "ProductUsage" ("Product_Id");
        CREATE INDEX IF NOT EXISTS "IX_ProductUsage_Researcher_Id" ON "ProductUsage" ("Researcher_Id");

        -- Insert into __EFMigrationsHistory
        INSERT INTO "__EFMigrationsHistory" ("MigrationId", "ProductVersion")
        VALUES ('20221022152249_InitialMigration', '6.0.5');
    END IF;
END
$EF$;
