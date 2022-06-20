/// Describes a write that creates a file.
pub(super) struct CreateFileAction {
    pub(crate) relative_path: String,
    pub(crate) content: Vec<u8>,
}

/// Describes a write that creates a directory.
pub(super) struct CreateDirectoryAction {
    pub(crate) relative_path: String,
}

/// A WriteAction is a description of a write that must happen sometime.
///
/// The reason this exists is twofold:
///
/// 1. Traits can't have async functions right now, and even if they do, dealing
///    with a lot of async functions is a pain. This affects us because the flow
///    is implemented in structs that implement traits, and if the flow does the
///    writes itself it'll have to be an async function.
/// 2. It's more limiting: I don't think the flows should be able to do writes
///    themselves, as they don't *need* to. I prefer the functional approach
///    of making as many functions as pure as possible, and performing all side
///    effects in as few places as possible.
pub(super) enum WriteAction {
    CreateFile(CreateFileAction),
    CreateDirectory(CreateDirectoryAction),
}

impl WriteAction {
    /// Asynchronously performs the write described by this action.
    pub(super) async fn write(self, base_path: &std::path::PathBuf) -> std::io::Result<()> {
        match self {
            WriteAction::CreateFile(file) => {
                let final_path = base_path.join(file.relative_path);
                tokio::fs::write(final_path, file.content).await
            }
            WriteAction::CreateDirectory(directory) => {
                let final_path = base_path.join(directory.relative_path);
                tokio::fs::create_dir_all(final_path).await
            }
        }
    }
}
