/* --------------------------------------------------------------------
// CreateDatabase.sql -- Creates the database for the part numbering sys. 
//
// Created 2018-09-20 DLB
// --------------------------------------------------------------------
*/

drop Database PnSysData;

Create Database PnSysData;
Use PnSysData;

create table Designers 
(
    Name   varchar(120),         /* Name of designer in F. LastName form */
    Year0  varchar(120),         /* Season that this designer was installed */
    Active int                   /* Non zero if active */
);

create table Projects
(
    ProjectId   char(3) PRIMARY KEY,
    Description varchar(512),
    Year0       varchar(10),      /* Year created, such as "Fall18" */
    Active      int               /* Non zero if active */
);

create table Subsystems
(
    ProjectId   char(3),          /* Must be in Projects Table */
    SubsystemId char(2),      
    Description varchar(512)
);

create table PartTypes 
(
    Digit       char(1),
    Description varchar (256)
);

create table EpicParts
(
    PID         char(32),       /* Unique internal ID assigned to this part */
    ProjectId   char(3),        /* Project, must be one of the above */
    SubsystemId char(2),        /* Subsystem, must be in sub system table */
    PartType    char(1),        /* Must be in the part-type table */
    SequenceNum int,            /* Incremented in golang code */
    Designer    varchar(60),            
    DateIssued  char(20),
    Description varchar(512)
);

create table SupplierCategory
(
    Category char(2),
    Description char(120)
);

create table SupplierParts
(
    PID         char(32),
    Category    char(2),
    SeqNum      int,
    Description varchar(512),
    Vendor      varchar(256),
    VendorPN    varchar(256),    
    WebLink     varchar(512),
    Designer    varchar(60),
    DateIssued  char(20)
);


insert into PartTypes (Digit, Description) values("0", "Assembly");
insert into PartTypes (Digit, Description) values("1", "Sheet Metal");
insert into PartTypes (Digit, Description) values("2", "Extrusion");
insert into PartTypes (Digit, Description) values("3", "3D Printed");
insert into PartTypes (Digit, Description) values("4", "Machined");
insert into PartTypes (Digit, Description) values("5", "Wood");
insert into PartTypes (Digit, Description) values("6", "Prototype");
insert into PartTypes (Digit, Description) values("7", "Volume");
insert into PartTypes (Digit, Description) values("8", "CNC");
insert into PartTypes (Digit, Description) values("9", "Misc");

insert into SupplierCategory (Category, Description) values("01", "Wheels");
insert into SupplierCategory (Category, Description) values("02", "Motors and Servos");
insert into SupplierCategory (Category, Description) values("03", "Gearboxes/Transmissions");
insert into SupplierCategory (Category, Description) values("04", "Drive Chassis");
insert into SupplierCategory (Category, Description) values("05", "Electrical Components");
insert into SupplierCategory (Category, Description) values("06", "Gears");
insert into SupplierCategory (Category, Description) values("07", "Bearings");
insert into SupplierCategory (Category, Description) values("08", "Hubs");
insert into SupplierCategory (Category, Description) values("09", "Sprockets");
insert into SupplierCategory (Category, Description) values("10", "Pneumatics");
insert into SupplierCategory (Category, Description) values("11", "Pulleys");
insert into SupplierCategory (Category, Description) values("12", "Game Pieces");
insert into SupplierCategory (Category, Description) values("13", "Screws, Nuts, and Bolts");
insert into SupplierCategory (Category, Description) values("14", "Extrusions");
insert into SupplierCategory (Category, Description) values("15", "Brackets, Tees, Connectors");
insert into SupplierCategory (Category, Description) values("16", "Washers and Spacers");
insert into SupplierCategory (Category, Description) values("17", "Collars");
insert into SupplierCategory (Category, Description) values("18", "Hinges");
insert into SupplierCategory (Category, Description) values("99", "Misc.");

/*
create table Pictures
(
    PicID char(32) NOT NULL UNIQUE,
    OriginalFileName varchar(512),
    Width            int,         /* Original Width //
    Height           int,         /* Original Height //
    FileSize         int,         /* Size in bytes, 0=not known //
    Extention        varchar(32), /* Picture extension with period //
    DateOfUpload     datetime
);

create table Users
(
   UserID       char(32) PRIMARY KEY,
   UserName     varchar(40) NOT NULL,  /* Cannot be changed! //
   PasswordHash varchar(60),                
   LastName     varchar(120),
   FirstName    varchar(120),
   NickName     varchar(120),
   Title        varchar(80),
   Email        varchar(512),
   Tags         varchar(512),      
   PicId        char(32),                 /* Zero if no pic //                  
   Active       boolean,
   DateCreated  datetime
);
*/