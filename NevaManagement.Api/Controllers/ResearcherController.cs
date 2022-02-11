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

        [HttpGet]
        public async Task<IActionResult> GetAll()
        {
            var researchers = await this.repository.GetAll();

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
