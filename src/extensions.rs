pub trait JoinableIterator<F>
where
    F: std::future::Future,
{
    fn join_all(self) -> futures::future::JoinAll<F>;
}

impl<T, Fu> JoinableIterator<Fu> for T
where
    T: std::iter::Iterator<Item = Fu>,
    Fu: std::future::Future,
{
    fn join_all(self) -> futures::future::JoinAll<Fu> {
        futures::future::join_all(self)
    }
}
