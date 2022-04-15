using NevaManagement.Domain.Dtos.Organism;
using NevaManagement.Domain.Interfaces.Repositories;
using NevaManagement.Domain.Interfaces.Services;
using NevaManagement.Domain.Models;
using System;
using System.Collections.Generic;
using System.Threading.Tasks;

namespace NevaManagement.Infrastructure.Services
{
    public class OrganismService : IOrganismService
    {
        private readonly IOrganismRepository repository;

        public OrganismService(IOrganismRepository repository)
        {
            this.repository = repository;
        }

        public async Task<IList<GetOrganismDto>> GetOrganisms()
        {
            return await this.repository.GetOrganisms();
        }

        public async Task<bool> AddOrganism(AddOrganismDto addOrganismDto)
        {
            var result = false;
            var organism = new Organism
            {
                Description = addOrganismDto.Description,
                Name = addOrganismDto.Name,
                CollectionDate = addOrganismDto.CollectionDate,
                CollectionLocation = addOrganismDto.CollectionLocation,
                IsolationDate = addOrganismDto.IsolationDate,
                OriginPart = addOrganismDto.OriginPart,
                Type = addOrganismDto.Type
            };

            if (addOrganismDto.OriginOrganismId is not null)
            {
                var originOrganism = await this.repository.GetEntityById(addOrganismDto.OriginOrganismId.Value);
                organism.Origin = originOrganism;
            }

            try
            {
                result = await this.repository.Create(organism);
            }
            catch
            {
                throw;
            }

            return result;
        }

        public async Task<bool> EditOrganism(EditOrganismDto editOrganismDto)
        {
            var result = false;
            var organism = await this.repository.GetEntityById(editOrganismDto.Id);

            if (organism is not null)
            {
                organism.Name = editOrganismDto.Name;
                organism.Description = editOrganismDto.Description;

                try
                {
                    result = await this.repository.SaveChanges();
                }
                catch
                {
                    throw;
                }

            }

            return result;
        }

        public async Task<GetDetailedOrganismDto> GetOrganismById(long organismId)
        {
            return await this.repository.GetById(organismId);
        }
    }
}
