/* --------------------------------------------------------------------
// LoadDatabase.sql -- Loads initial data for part systems. 
//
// Created 2018-09-24 DLB
// --------------------------------------------------------------------
*/

Use PnSysData;

insert into Designers (Name, Year0, Active) values("U. Unknown",    "2010", 0);
insert into Designers (Name, Year0, Active) values("A. Yanney",     "2017", 1);
insert into Designers (Name, Year0, Active) values("B. Long",       "2014", 1);
insert into Designers (Name, Year0, Active) values("D. Brandon",    "2013", 1);
insert into Designers (Name, Year0, Active) values("D. Fieldhouse", "2013", 1);
insert into Designers (Name, Year0, Active) values("H. Fieldhouse", "2015", 1);
insert into Designers (Name, Year0, Active) values("J. Bassett",    "2015", 1);
insert into Designers (Name, Year0, Active) values("J. Stratton",   "2017", 1);
insert into Designers (Name, Year0, Active) values("J. Wu",         "2016", 0);
insert into Designers (Name, Year0, Active) values("K. Dominguez",  "2012", 0);
insert into Designers (Name, Year0, Active) values("M. Schroeder",  "2015", 1);
insert into Designers (Name, Year0, Active) values("P. DeVries",    "2012", 1);
insert into Designers (Name, Year0, Active) values("R. Carranza",   "2015", 0);
insert into Designers (Name, Year0, Active) values("R. Rietveld",   "2015", 1);
insert into Designers (Name, Year0, Active) values("R. Vreeke",     "2014", 0);
insert into Designers (Name, Year0, Active) values("S. Tefera",     "2014", 0);
insert into Designers (Name, Year0, Active) values("S. Wang",       "2014", 0);
insert into Designers (Name, Year0, Active) values("T. Qu",         "2016", 1);
insert into Designers (Name, Year0, Active) values("C. Dominik",    "2013", 0);
insert into Designers (Name, Year0, Active) values("C. Granch",     "2014", 0);
insert into Designers (Name, Year0, Active) values("L. Furlong",    "2014", 0);
insert into Designers (Name, Year0, Active) values("N. Gardner",    "2015", 0);
insert into Designers (Name, Year0, Active) values("N. Dekker",     "2014", 0);
insert into Designers (Name, Year0, Active) values("Y. Ghebreyesus","2015", 0);
insert into Designers (Name, Year0, Active) values("C. Yuan",       "2014", 0);
insert into Designers (Name, Year0, Active) values("A. Haitz",      "2014", 0);
insert into Designers (Name, Year0, Active) values("D. OToole",     "2014", 0);

insert into Projects (ProjectId, Description, Year0, Active) values("R19", "Competition Prototype 2019",  "2019", 0);
insert into Projects (ProjectId, Description, Year0, Active) values("C19", "Competition Final 2019",      "2019", 0);
insert into Projects (ProjectId, Description, Year0, Active) values("G18", "General 2018",                "2018", 1);
insert into Projects (ProjectId, Description, Year0, Active) values("E18", "Practice Bot Fall 2018",      "2018", 1);
insert into Projects (ProjectId, Description, Year0, Active) values("R18", "Competition Prototype 2018",  "2018", 0);
insert into Projects (ProjectId, Description, Year0, Active) values("C18", "Competition Final 2018",      "2018", 0);
insert into Projects (ProjectId, Description, Year0, Active) values("N18", "Conceptual Bot 2018",         "2018", 0);
insert into Projects (ProjectId, Description, Year0, Active) values("F18", "Field Elements 2018",         "2018", 0);
insert into Projects (ProjectId, Description, Year0, Active) values("G17", "General (Fall 2017)",         "2017", 0);
insert into Projects (ProjectId, Description, Year0, Active) values("P17", "Pre-Season (Fall 2017)",      "2018", 0);
insert into Projects (ProjectId, Description, Year0, Active) values("V18", "Volumentric Bot 2018",        "2018", 0);
insert into Projects (ProjectId, Description, Year0, Active) values("FG8", "Frisbee Go Kart (Summer 18)", "2018", 0);
insert into Projects (ProjectId, Description, Year0, Active) values("P16", "Preseason (Fall 2016)",       "2016", 0);
insert into Projects (ProjectId, Description, Year0, Active) values("S17", "Season 2017",                 "2017", 0);
insert into Projects (ProjectId, Description, Year0, Active) values("C17", "Competition Bot (2017)",      "2017", 0);
insert into Projects (ProjectId, Description, Year0, Active) values("B17", "Business Projects (2017)",    "2017", 0);
insert into Projects (ProjectId, Description, Year0, Active) values("F17", "Field Support Equpipment",    "2017", 0);
insert into Projects (ProjectId, Description, Year0, Active) values("T17", "Post-season (Spring 2017)",   "2017", 0);
insert into Projects (ProjectId, Description, Year0, Active) values("SB5", "SW Class Beg (Spring 2015)",  "2015", 0);
insert into Projects (ProjectId, Description, Year0, Active) values("SM5", "SW Class Int (Spring 2015)",  "2015", 0);
insert into Projects (ProjectId, Description, Year0, Active) values("P15", "Preseason 2015",              "2015", 0);
insert into Projects (ProjectId, Description, Year0, Active) values("S16", "Season 2016",                 "2016", 0);
insert into Projects (ProjectId, Description, Year0, Active) values("S14", "SoildWorks Class Fall 2014",  "2014", 0);
insert into Projects (ProjectId, Description, Year0, Active) values("F15", "FRC 2015 Competition Bot",    "2015", 0);
insert into Projects (ProjectId, Description, Year0, Active) values("P15", "Prototype 2015",              "2015", 0);
insert into Projects (ProjectId, Description, Year0, Active) values("C15", "Competition 2015",            "2015", 0);
insert into Projects (ProjectId, Description, Year0, Active) values("SPP", "Special Projects 2014",       "2014", 0);
insert into Projects (ProjectId, Description, Year0, Active) values("Y14", "All of Year 2014 Projects",   "2014", 0);

insert into Subsystems (SubsystemId, ProjectId, Description) values("GN", "G18", "General");
insert into Subsystems (SubsystemId, ProjectId, Description) values("CH", "F18", "Chassis");
insert into Subsystems (SubsystemId, ProjectId, Description) values("FS", "F18", "Frisbee Shooter");
insert into Subsystems (SubsystemId, ProjectId, Description) values("AR", "C18", "Arm");
insert into Subsystems (SubsystemId, ProjectId, Description) values("BF", "C18", "Butterfly Drive - Cassette Design");
insert into Subsystems (SubsystemId, ProjectId, Description) values("CH", "C18", "Chassis");
insert into Subsystems (SubsystemId, ProjectId, Description) values("CL", "C18", "Climber");
insert into Subsystems (SubsystemId, ProjectId, Description) values("ER", "C18", "Entire Robot");
insert into Subsystems (SubsystemId, ProjectId, Description) values("LW", "C18", "Power Cube Grabber");
insert into Subsystems (SubsystemId, ProjectId, Description) values("TL", "C18", "Telescoping Tower of Terror");
insert into Subsystems (SubsystemId, ProjectId, Description) values("FE", "F18", "Field Elements");
insert into Subsystems (SubsystemId, ProjectId, Description) values("GN", "G17", "General");
insert into Subsystems (SubsystemId, ProjectId, Description) values("BB", "N18", "Lidar Turntable (Bobby Box)");
insert into Subsystems (SubsystemId, ProjectId, Description) values("BF", "N18", "Butterfly Drive - Cassette Design");
insert into Subsystems (SubsystemId, ProjectId, Description) values("CH", "N18", "Chassis - Conceptual");
insert into Subsystems (SubsystemId, ProjectId, Description) values("CL", "N18", "Climber - Conceputal");
insert into Subsystems (SubsystemId, ProjectId, Description) values("ER", "N18", "Entire Robot");
insert into Subsystems (SubsystemId, ProjectId, Description) values("FK", "N18", "Forklift Design");
insert into Subsystems (SubsystemId, ProjectId, Description) values("LD", "N18", "Ladder Arm Design");
insert into Subsystems (SubsystemId, ProjectId, Description) values("LW", "N18", "Claw Cube Intake");
insert into Subsystems (SubsystemId, ProjectId, Description) values("MS", "N18", "Mini Folding Arm (The Mini Sam)");
insert into Subsystems (SubsystemId, ProjectId, Description) values("TL", "N18", "Telescoping Arm");
insert into Subsystems (SubsystemId, ProjectId, Description) values("WD", "N18", "Wedge Ladder Design");
insert into Subsystems (SubsystemId, ProjectId, Description) values("WH", "N18", "Wheel Cube Intake");
insert into Subsystems (SubsystemId, ProjectId, Description) values("CH", "P17", "Chassis (6 Wheels)");
insert into Subsystems (SubsystemId, ProjectId, Description) values("CS", "P17", "CNC Support");
insert into Subsystems (SubsystemId, ProjectId, Description) values("FL", "P17", "FLL Project");
insert into Subsystems (SubsystemId, ProjectId, Description) values("SW", "P17", "Ball Sweeper");
insert into Subsystems (SubsystemId, ProjectId, Description) values("TS", "P17", "Targeting System");
insert into Subsystems (SubsystemId, ProjectId, Description) values("CH", "R18", "Chassis");
insert into Subsystems (SubsystemId, ProjectId, Description) values("CH", "V18", "Chassis");
insert into Subsystems (SubsystemId, ProjectId, Description) values("CL", "V18", "Climber");
insert into Subsystems (SubsystemId, ProjectId, Description) values("FA", "V18", "Folding Arm Design");
insert into Subsystems (SubsystemId, ProjectId, Description) values("GN", "V18", "General");
insert into Subsystems (SubsystemId, ProjectId, Description) values("LD", "V18", "Ladder Arm Design");
insert into Subsystems (SubsystemId, ProjectId, Description) values("LW", "V18", "Claw Cube Griper");
insert into Subsystems (SubsystemId, ProjectId, Description) values("PD", "V18", "Pendulum Arm Design");
insert into Subsystems (SubsystemId, ProjectId, Description) values("EE", "C18", "Electronics");
insert into Subsystems (SubsystemId, ProjectId, Description) values("FR", "FG8", "Go Kart - Frame");
insert into Subsystems (SubsystemId, ProjectId, Description) values("SS", "FG8", "Go Kart - Steering System");
insert into Subsystems (SubsystemId, ProjectId, Description) values("DS", "FG8", "Go Kart - Driveshaft");
insert into Subsystems (SubsystemId, ProjectId, Description) values("ST", "FG8", "Go Kart - Seat");
insert into Subsystems (SubsystemId, ProjectId, Description) values("PD", "FG8", "Go Kart - Pedals");
insert into Subsystems (SubsystemId, ProjectId, Description) values("MH", "P16", "Marble Shooter");
insert into Subsystems (SubsystemId, ProjectId, Description) values("DS", "P16", "Drivers Station");
insert into Subsystems (SubsystemId, ProjectId, Description) values("BH", "P16", "Ball Shooter");
insert into Subsystems (SubsystemId, ProjectId, Description) values("PB", "P16", "Protobot");
insert into Subsystems (SubsystemId, ProjectId, Description) values("MC", "P16", "Motor Controller");
insert into Subsystems (SubsystemId, ProjectId, Description) values("SG", "P16", "EPIC Board Sign");
insert into Subsystems (SubsystemId, ProjectId, Description) values("CH", "S17", "Chassis");
insert into Subsystems (SubsystemId, ProjectId, Description) values("VI", "S17", "Volume Intergration");
insert into Subsystems (SubsystemId, ProjectId, Description) values("GN", "S17", "General (Season 2017)");
insert into Subsystems (SubsystemId, ProjectId, Description) values("BO", "S17", "Ball Shooter - Single Wheel");
insert into Subsystems (SubsystemId, ProjectId, Description) values("BD", "S17", "Ball Shooter - Dual Wheel");
insert into Subsystems (SubsystemId, ProjectId, Description) values("CL", "S17", "Climber");
insert into Subsystems (SubsystemId, ProjectId, Description) values("GH", "S17", "Gear Holder");
insert into Subsystems (SubsystemId, ProjectId, Description) values("FT", "S17", "Fuel Tank");
insert into Subsystems (SubsystemId, ProjectId, Description) values("FS", "S17", "Fuel Shooter");
insert into Subsystems (SubsystemId, ProjectId, Description) values("FC", "S17", "Fuel Collector");
insert into Subsystems (SubsystemId, ProjectId, Description) values("SB", "S17", "Ball Shooter - Softball");
insert into Subsystems (SubsystemId, ProjectId, Description) values("ER", "S17", "Entire Robot");
insert into Subsystems (SubsystemId, ProjectId, Description) values("FL", "S17", "Fuel Loader");
insert into Subsystems (SubsystemId, ProjectId, Description) values("SL", "S17", "Shooter Loader");
insert into Subsystems (SubsystemId, ProjectId, Description) values("CL", "C17", "Climber (Competition)");
insert into Subsystems (SubsystemId, ProjectId, Description) values("BO", "C17", "Softball Shooter (Competition)");
insert into Subsystems (SubsystemId, ProjectId, Description) values("ER", "C17", "Entire Robot (Competition)");
insert into Subsystems (SubsystemId, ProjectId, Description) values("FT", "C17", "Fuel Tank (Competition)");
insert into Subsystems (SubsystemId, ProjectId, Description) values("FC", "C17", "Fuel Collector (Competition)");
insert into Subsystems (SubsystemId, ProjectId, Description) values("FL", "C17", "Fuel Loader (Competition)");
insert into Subsystems (SubsystemId, ProjectId, Description) values("TT", "C17", "Twin Turret (Competition)");
insert into Subsystems (SubsystemId, ProjectId, Description) values("SL", "C17", "Shooter Loader (Competition)");
insert into Subsystems (SubsystemId, ProjectId, Description) values("CH", "C17", "Chassis (Competition)");
insert into Subsystems (SubsystemId, ProjectId, Description) values("GH", "C17", "Gear Holder (Competition)");
insert into Subsystems (SubsystemId, ProjectId, Description) values("GG", "C17", "Gear Grabber (Competition)");
insert into Subsystems (SubsystemId, ProjectId, Description) values("BL", "C17", "Blender (Competition)");
insert into Subsystems (SubsystemId, ProjectId, Description) values("SV", "B17", "Pit Shelves");
insert into Subsystems (SubsystemId, ProjectId, Description) values("CT", "F17", "Bot Cart");
insert into Subsystems (SubsystemId, ProjectId, Description) values("SP", "C17", "Spear");
insert into Subsystems (SubsystemId, ProjectId, Description) values("DS", "F17", "Drivers Station");
insert into Subsystems (SubsystemId, ProjectId, Description) values("BG", "T17", "Beginner CAD, SP 2017");
insert into Subsystems (SubsystemId, ProjectId, Description) values("ID", "T17", "Intermediate CAD, SP 2017");
insert into Subsystems (SubsystemId, ProjectId, Description) values("AV", "T17", "Advanced CAD, SP 2017");
insert into Subsystems (SubsystemId, ProjectId, Description) values("DE", "SB5", "Demo Parts");
insert into Subsystems (SubsystemId, ProjectId, Description) values("DE", "SM5", "Demo Parts");
insert into Subsystems (SubsystemId, ProjectId, Description) values("GE", "SM5", "General");
insert into Subsystems (SubsystemId, ProjectId, Description) values("FW", "P15", "Flywheel");
insert into Subsystems (SubsystemId, ProjectId, Description) values("SP", "P15", "Special Projects");
insert into Subsystems (SubsystemId, ProjectId, Description) values("CN", "P15", "CNC");
insert into Subsystems (SubsystemId, ProjectId, Description) values("MF", "P15", "Mini Frissbee");
insert into Subsystems (SubsystemId, ProjectId, Description) values("CD", "S16", "Conceptual Designs");
insert into Subsystems (SubsystemId, ProjectId, Description) values("PR", "P15", "Prototyping Bots");
insert into Subsystems (SubsystemId, ProjectId, Description) values("WH", "S16", "Wheels");
insert into Subsystems (SubsystemId, ProjectId, Description) values("PN", "S16", "Pneumatics");
insert into Subsystems (SubsystemId, ProjectId, Description) values("CH", "S16", "Chassis");
insert into Subsystems (SubsystemId, ProjectId, Description) values("WG", "S16", "Defense Manipulator (Wedgie)");
insert into Subsystems (SubsystemId, ProjectId, Description) values("SH", "S16", "Ball Shooter");
insert into Subsystems (SubsystemId, ProjectId, Description) values("CL", "S16", "Climber");
insert into Subsystems (SubsystemId, ProjectId, Description) values("BB", "S16", "Bailer Bar");
insert into Subsystems (SubsystemId, ProjectId, Description) values("ER", "S16", "Entire Robot Assy");
insert into Subsystems (SubsystemId, ProjectId, Description) values("DS", "S16", "Drivers Station");
insert into Subsystems (SubsystemId, ProjectId, Description) values("EL", "S16", "Electronics");
insert into Subsystems (SubsystemId, ProjectId, Description) values("CH", "F15", "Chassis");
insert into Subsystems (SubsystemId, ProjectId, Description) values("GA", "F15", "General");
insert into Subsystems (SubsystemId, ProjectId, Description) values("AR", "S14", "Arm");
insert into Subsystems (SubsystemId, ProjectId, Description) values("BG", "S14", "Assy Project");
insert into Subsystems (SubsystemId, ProjectId, Description) values("KJ", "P15", "Kyles");
insert into Subsystems (SubsystemId, ProjectId, Description) values("CC", "C15", "Caleb claw");
insert into Subsystems (SubsystemId, ProjectId, Description) values("HG", "C15", "Happy Gilmore");
insert into Subsystems (SubsystemId, ProjectId, Description) values("CH", "C15", "Chassis 2015");
insert into Subsystems (SubsystemId, ProjectId, Description) values("FW", "C15", "Folding wing");
insert into Subsystems (SubsystemId, ProjectId, Description) values("EL", "C15", "Electronics");
insert into Subsystems (SubsystemId, ProjectId, Description) values("EV", "C15", "Elevator");
insert into Subsystems (SubsystemId, ProjectId, Description) values("GR", "C15", "General Robot Assy");
insert into Subsystems (SubsystemId, ProjectId, Description) values("PT", "SPP", "Pit");
insert into Subsystems (SubsystemId, ProjectId, Description) values("GA", "C15", "General");

insert into Subsystems (SubsystemId, ProjectId, Description) values("S2", "Y14", "Old subsystem.");
insert into Subsystems (SubsystemId, ProjectId, Description) values("A1", "Y14", "Old subsystem.");
insert into Subsystems (SubsystemId, ProjectId, Description) values("M1", "Y14", "Old subsystem.");
insert into Subsystems (SubsystemId, ProjectId, Description) values("L1", "Y14", "Old subsystem.");
insert into Subsystems (SubsystemId, ProjectId, Description) values("FL", "Y14", "Old subsystem.");
insert into Subsystems (SubsystemId, ProjectId, Description) values("CH", "Y14", "Old subsystem.");
insert into Subsystems (SubsystemId, ProjectId, Description) values("B1", "Y14", "Old subsystem.");
insert into Subsystems (SubsystemId, ProjectId, Description) values("AM", "Y14", "Old subsystem.");
insert into Subsystems (SubsystemId, ProjectId, Description) values("MM", "Y14", "Old subsystem.");
insert into Subsystems (SubsystemId, ProjectId, Description) values("L3", "Y14", "Old subsystem.");

