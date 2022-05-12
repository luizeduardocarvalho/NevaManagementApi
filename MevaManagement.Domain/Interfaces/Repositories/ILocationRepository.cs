﻿namespace NevaManagement.Domain.Interfaces.Repositories;

public interface ILocationRepository
{
    Task<GetDetailedLocationDto> GetById(long id);
    Task<Location> GetEntityById(long id);
    Task<IList<GetLocationDto>> GetLocations();
    Task<bool> Create(Location location);
    Task<bool> SaveChanges();
}
