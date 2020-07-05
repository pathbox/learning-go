Should I use it?

Well, probably not. Don’t be fooled by the simple interface. As I mentioned before, the atomic package is not for most people, and even an innocent looking feature like this one has some rough edges. Quoting the documentation:

 All calls to Store for a given Value must use values of the same concrete type. Store of an inconsistent type panics, as does Store(nil).
or

 Once Store has been called, a Value must not be copied.
Now that you just read about it you know it, but someone else in your team might not, or you might not remember it 6 months from now. This is why you shouldn’t use this unless you really need to worry about contention or you are operating at a scale where this performance improvements make a difference.

So, which kind to applications am I talking about? you may ask. Why is this slightly tricky thing there in the first place? It turns out that this benchmark in my modest laptop was a bit out of scope