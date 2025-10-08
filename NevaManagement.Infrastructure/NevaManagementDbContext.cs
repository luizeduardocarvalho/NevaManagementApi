namespace NevaManagement.Infrastructure;

public class NevaManagementDbContext : DbContext
{
    public DbSet<Article> Articles { get; set; }
    public DbSet<ArticleContainer> ArticleContainers { get; set; }
    public DbSet<Container> Containers { get; set; }
    public DbSet<Equipment> Equipments { get; set; }
    public DbSet<EquipmentUsage> EquipmentUsages { get; set; }
    public DbSet<Location> Locations { get; set; }
    public DbSet<Organism> Organisms { get; set; }
    public DbSet<Product> Products { get; set; }
    public DbSet<ProductUsage> ProductUsages { get; set; }
    public DbSet<Researcher> Researchers { get; set; }
    public DbSet<Laboratory> Laboratories { get; set; }
    public DbSet<LaboratoryInvitation> LaboratoryInvitations { get; set; }

    public NevaManagementDbContext(DbContextOptions<NevaManagementDbContext> options)
        : base(options)
    {
    }
}