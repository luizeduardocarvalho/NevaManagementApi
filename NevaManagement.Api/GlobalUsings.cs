﻿global using System;
global using System.Collections.Generic;
global using System.IdentityModel.Tokens.Jwt;
global using System.Net;
global using System.Security.Claims;
global using System.Text;
global using System.Text.Json;
global using System.Threading.Tasks;
global using AutoMapper;
global using Microsoft.AspNetCore.Authentication.JwtBearer;
global using Microsoft.AspNetCore.Authorization;
global using Microsoft.AspNetCore.Builder;
global using Microsoft.AspNetCore.Hosting;
global using Microsoft.AspNetCore.Http;
global using Microsoft.AspNetCore.Mvc;
global using Microsoft.EntityFrameworkCore;
global using Microsoft.Extensions.Configuration;
global using Microsoft.Extensions.DependencyInjection;
global using Microsoft.Extensions.Hosting;
global using Microsoft.Extensions.Options;
global using Microsoft.IdentityModel.Tokens;
global using Microsoft.OpenApi.Models;
global using NevaManagement.Api.Configurations;
global using NevaManagement.Api.Middlewares;
global using NevaManagement.Domain.Dtos.Auth;
global using NevaManagement.Domain.Dtos.Container;
global using NevaManagement.Domain.Dtos.Equipment;
global using NevaManagement.Domain.Dtos.EquipmentUsage;
global using NevaManagement.Domain.Dtos.Organism;
global using NevaManagement.Domain.Dtos.Product;
global using NevaManagement.Domain.Dtos.Researcher;
global using NevaManagement.Domain.Interfaces.Repositories;
global using NevaManagement.Domain.Interfaces.Services;
global using NevaManagement.Domain.Models;
global using NevaManagement.Infrastructure;
global using NevaManagement.Infrastructure.Repositories;
global using NevaManagement.Infrastructure.Services;
