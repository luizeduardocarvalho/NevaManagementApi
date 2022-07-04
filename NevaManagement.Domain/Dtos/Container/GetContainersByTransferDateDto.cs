namespace NevaManagement.Domain.Dtos.Container;

public class GetContainersByTransferDateDto
{
    public long Id { get; set; }

    public DateTimeOffset TransferDate { get; set; }

    public string Name { get; set; }
}
