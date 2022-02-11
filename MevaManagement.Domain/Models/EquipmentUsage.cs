using System;
using System.ComponentModel.DataAnnotations.Schema;

namespace NevaManagement.Domain.Models
{
    [Table("EquipmentUsage")]
    public class EquipmentUsage : BaseEntity
    {
        [ForeignKey("Researcher_Id")]
        public Researcher Researcher { get; set; }

        [ForeignKey("Equipment_Id")]
        public Equipment Equipment { get; set; }

        public DateTimeOffset UsageDate { get; set; }

        public string Description { get; set; }
    }
}
