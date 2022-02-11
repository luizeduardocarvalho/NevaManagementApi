using System;
using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;

namespace NevaManagement.Domain.Models
{
    [Table("Container")]
    public class Container : BaseEntity
    {
        [MaxLength(80)]
        public string Name { get; set; }

        public string Description { get; set; }

        [Column("Culture_Media")]
        public string CultureMedia { get; set; }

        [ForeignKey("Equipment_Id")]
        public Equipment Equipment { get; set; }

        [ForeignKey("Researcher_Id")]
        public Researcher Researcher { get; set; }

        [ForeignKey("Origin_Id")]
        public Organism Origin { get; set; }

        [ForeignKey("Sub_Container_Id")]
        public Container SubContainer { get; set; }

        [Column("Creation_Date")]
        public DateTimeOffset CreationDate { get; set; }
    }
}
