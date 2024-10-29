#include "mesh_conn_local_rx.h"
#include <string.h>

namespace mesh {

namespace connection {

using namespace mesh::connection;

LocalRx::LocalRx(memif_ops_t *ops) : Local()
{
    _kind = Kind::receiver;

    memcpy(&this->ops, ops, sizeof(memif_ops_t));
}

LocalRx::~LocalRx()
{
}

Result LocalRx::on_establish(context::Context& ctx)
{
    return Result::success;
}

Result LocalRx::on_shutdown(context::Context& ctx)
{
    return Result::success;
}

void LocalRx::on_delete(context::Context& ctx)
{
}

} // namespace connection

} // namespace mesh
