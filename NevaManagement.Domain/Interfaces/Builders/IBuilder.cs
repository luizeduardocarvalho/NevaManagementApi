namespace NevaManagement.Domain.Interfaces.Builders;

public interface IBuilder<in TIn, out TOut> 
    where TIn : notnull
{
    TOut Build(TIn input);
}
