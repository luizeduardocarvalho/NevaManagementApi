#See https://aka.ms/containerfastmode to understand how Visual Studio uses this Dockerfile to build your images for faster debugging.

FROM mcr.microsoft.com/dotnet/aspnet:5.0 AS base
WORKDIR /app
EXPOSE 80
EXPOSE 443

FROM mcr.microsoft.com/dotnet/sdk:5.0 AS build
WORKDIR /src
COPY ["NevaManagement.Api/NevaManagement.Api.csproj", "NevaManagement.Api/"]
COPY ["NevaManagement.Infrastructure/NevaManagement.Infrastructure.csproj", "NevaManagement.Infrastructure/"]
COPY ["MevaManagement.Domain/NevaManagement.Domain.csproj", "MevaManagement.Domain/"]
RUN dotnet restore "NevaManagement.Api/NevaManagement.Api.csproj"
COPY . .
WORKDIR "/src/NevaManagement.Api"
RUN dotnet build "NevaManagement.Api.csproj" -c Release -o /app/build

FROM build AS publish
RUN dotnet publish "NevaManagement.Api.csproj" -c Release -o /app/publish

FROM base AS final
WORKDIR /app
COPY --from=publish /app/publish .
CMD ASPNETCORE_URLS=http://*:$PORT dotnet NevaManagement.Api.dll