#include "mesh_conn_local_tx.h"

namespace mesh {

namespace connection {

using namespace mesh::connection;

LocalTx::LocalTx() : Local()
{
    _kind = Kind::transmitter;
}

LocalTx::~LocalTx()
{
}

Result LocalTx::on_establish(context::Context& ctx)
{
    // TODO: add implementation

    set_state(ctx, State::active);

    return Result::success;
}

Result LocalTx::on_receive(context::Context& ctx, void *ptr, uint32_t sz,
                           uint32_t& sent)
{
    // TODO: add implementation

    std::string str((char *)ptr, sz);

    printf("Received: %s\n", str.c_str());

    return Result::success;
}

Result LocalTx::on_shutdown(context::Context& ctx)
{
    // TODO: add implementation

    set_state(ctx, State::closed);

    return Result::success;
}

void LocalTx::on_delete(context::Context& ctx)
{
}

} // namespace connection

} // namespace mesh
