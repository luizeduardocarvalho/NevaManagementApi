namespace NevaManagement.Infrastructure.Repositories;

public class BaseRepository<T> : IBaseRepository<T> where T : class
{
    private readonly NevaManagementDbContext context;
    public DbSet<T> table = null;

    public BaseRepository(NevaManagementDbContext context)
    {
        this.context = context;
        table = this.context.Set<T>();
    }

    public async Task Insert(T entity)
    {
        await this.table.AddAsync(entity);
    }

    public async Task<bool> SaveChanges()
    {
        var result = await this.context.SaveChangesAsync();

        if (result > 0)
        {
            return true;
        }

        return false;
    }
}
