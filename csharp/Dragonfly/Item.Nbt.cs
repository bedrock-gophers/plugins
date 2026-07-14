#nullable enable

namespace Dragonfly;

internal static class ItemNbtCodec
{
    internal static bool TryEncode(World.Item? item, out byte[] data)
    {
        switch (item)
        {
            case Item.BookAndQuill book:
                data = EncodeBookAndQuill(book);
                return true;
            case Item.WrittenBook book:
                data = EncodeWrittenBook(book);
                return true;
            default:
                data = [];
                return false;
        }
    }

    internal static World.Item Decode(World.Item item, ReadOnlySpan<byte> data, out bool consumed)
    {
        switch (item)
        {
            case Item.BookAndQuill book:
                consumed = true;
                return DecodeBookAndQuill(book, data);
            case Item.WrittenBook book:
                consumed = true;
                return DecodeWrittenBook(book, data);
            default:
                consumed = false;
                return item;
        }
    }

    private static byte[] EncodeBookAndQuill(Item.BookAndQuill book)
    {
        if (book.TotalPages() == 0) return [];
        return Nbt.Encode(new Nbt.Compound
        {
            ["pages"] = EncodePages(book.Pages),
        });
    }

    private static Item.BookAndQuill DecodeBookAndQuill(Item.BookAndQuill book, ReadOnlySpan<byte> data)
    {
        if (data.IsEmpty) return book;
        var pages = new List<string>(book.Pages);
        if (Nbt.Decode(data).TryGetValue("pages", out var encoded) && encoded.Type == Nbt.TagType.List)
        {
            foreach (var page in encoded.AsList())
            {
                if (page.Type != Nbt.TagType.Compound ||
                    !page.AsCompound().TryGetValue("text", out var text) || text.Type != Nbt.TagType.String)
                    continue;
                pages.Add(text.AsString());
            }
        }
        return new Item.BookAndQuill(pages.ToArray());
    }

    private static byte[] EncodeWrittenBook(Item.WrittenBook book) => Nbt.Encode(new Nbt.Compound
    {
        ["pages"] = EncodePages(book.Pages),
        ["author"] = Nbt.Value.String(book.Author),
        ["title"] = Nbt.Value.String(book.Title),
        ["generation"] = Nbt.Value.Byte(book.Generation.Uint8()),
    });

    private static Item.WrittenBook DecodeWrittenBook(Item.WrittenBook book, ReadOnlySpan<byte> data)
    {
        if (data.IsEmpty) return book;
        var root = Nbt.Decode(data);
        var pages = DecodePages(root);
        var title = GetString(root, "title");
        var author = GetString(root, "author");
        var generation = Item.OriginalGeneration();
        if (root.TryGetValue("generation", out var encodedGeneration) && encodedGeneration.Type == Nbt.TagType.Byte)
        {
            generation = encodedGeneration.AsByte() switch
            {
                1 => Item.CopyGeneration(),
                2 => Item.CopyOfCopyGeneration(),
                _ => Item.OriginalGeneration(),
            };
        }
        return new Item.WrittenBook(title, author, generation, pages);
    }

    private static Nbt.Value EncodePages(IEnumerable<string> pages) => Nbt.Value.List(
        Nbt.TagType.Compound,
        pages.Select(page => Nbt.Value.Compound(new Nbt.Compound
        {
            ["text"] = Nbt.Value.String(page),
        })).ToArray());

    private static string[] DecodePages(Nbt.Compound root)
    {
        if (!root.TryGetValue("pages", out var encoded) || encoded.Type != Nbt.TagType.List) return [];
        var pages = new List<string>();
        foreach (var page in encoded.AsList())
        {
            if (page.Type != Nbt.TagType.Compound ||
                !page.AsCompound().TryGetValue("text", out var text) || text.Type != Nbt.TagType.String)
                throw new InvalidDataException("invalid written book page");
            pages.Add(text.AsString());
        }
        return pages.ToArray();
    }

    private static string GetString(Nbt.Compound root, string name) =>
        root.TryGetValue(name, out var value) && value.Type == Nbt.TagType.String
            ? value.AsString()
            : string.Empty;
}
