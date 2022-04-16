using Microsoft.AspNetCore.Mvc;
using NevaManagement.Domain.Dtos.Researcher;
using NevaManagement.Domain.Interfaces.Repositories;
using NevaManagement.Domain.Models;
using System.Threading.Tasks;

namespace NevaManagement.Api.Controllers
{
    [ApiController]
    [Route("[controller]")]
    public class ResearcherController : ControllerBase
    {
        private readonly IResearcherRepository repository;

        public ResearcherController(IResearcherRepository repository)
        {
            this.repository = repository;
        }

        [HttpGet("GetResearchers")]
        public async Task<IActionResult> GetResearchers()
        {
            var researchers = await this.repository.GetResearchers();

            return Ok(researchers);
        }

        [HttpPost]
        public async Task<IActionResult> Create(CreateResearcherDto researcherDto)
        {
            if(!ModelState.IsValid)
            {
                return BadRequest();
            }

            var researcher = new Researcher()
            {
                Email = researcherDto.Email,
                Name = researcherDto.Name
            };

            var result = await this.repository.Create(researcher);

            if(result)
            {
                return Ok(result);
            }

            return StatusCode(500, "Error");
        }
    }
}
