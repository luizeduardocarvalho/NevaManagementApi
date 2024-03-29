﻿namespace NevaManagement.Infrastructure.Repositories;

public class BaseRepository<T> : IBaseRepository<T> where T : class
{
    private readonly NevaManagementDbContext context;
    public DbSet<T> table = null;

    public BaseRepository(NevaManagementDbContext context)
    {
        this.context = context;
        table = this.context.Set<T>();
    }

    public async Task<T> GetById(long id)
    {
        return await this.table.FindAsync(id);
    }

    public async Task Insert(T entity)
    {
        await this.table.AddAsync(entity);
    }

    public async Task<bool> SaveChanges()
    {
        var result = await this.context.SaveChangesAsync();

        return result > 0;
    }
}
