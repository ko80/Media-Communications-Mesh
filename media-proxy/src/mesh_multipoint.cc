#include "mesh_multipoint.h"

namespace mesh {

namespace connection {

using namespace mesh::connection;

MultipointGroup::MultipointGroup() : Connection()
{
    _kind = Kind::transmitter;
}

MultipointGroup::~MultipointGroup()
{
}

void MultipointGroup::register_out_link(Connection *link)
{
    out_link = link;
}

Result MultipointGroup::on_establish(context::Context& ctx)
{
    set_state(ctx, State::active);
    set_status(ctx, Status::healthy);

    return Result::success;
}

Result MultipointGroup::on_receive(context::Context& ctx, void *ptr,
                                   uint32_t sz, uint32_t& sent)
{
    if (!out_link)
        return Result::error_no_link_assigned;

    sent = sz;

    uint32_t out_sent = 0;

    out_link->do_receive(ctx, ptr, sz, out_sent);

    // DRAFT

    return Result::success;
}

Result MultipointGroup::on_shutdown(context::Context& ctx)
{
    set_state(ctx, State::closed);
    set_status(ctx, Status::shutdown);

    return Result::success;
}

void MultipointGroup::on_delete(context::Context& ctx)
{
}

} // namespace connection

} // namespace mesh
