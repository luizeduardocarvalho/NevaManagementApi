namespace NevaManagement.Domain.Interfaces.Repositories;

public interface IBaseRepository<T>
{
    Task<T> GetById(long id);
    Task<bool> SaveChanges();
    Task Insert(T entity);
}
