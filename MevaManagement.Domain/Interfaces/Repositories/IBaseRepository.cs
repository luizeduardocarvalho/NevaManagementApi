using System.Threading.Tasks;

namespace NevaManagement.Domain.Interfaces.Repositories
{
    public interface IBaseRepository<T>
    {
        Task<bool> SaveChanges();
        Task Insert(T entity);
    }
}
