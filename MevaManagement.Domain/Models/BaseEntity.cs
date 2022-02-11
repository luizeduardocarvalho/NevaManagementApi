using System.ComponentModel.DataAnnotations;

namespace NevaManagement.Domain.Models
{
    public class BaseEntity
    {
        [Key]
        public long Id { get; set; }
    }
}
