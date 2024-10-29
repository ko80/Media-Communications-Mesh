#ifndef MESH_MULTIPOINT_H
#define MESH_MULTIPOINT_H

#include "mesh_conn.h"

namespace mesh {

namespace connection {

class MultipointGroup : public Connection {

public:
    MultipointGroup();
    ~MultipointGroup() override;

    void register_out_link(Connection *link);

private:
    Result on_establish(context::Context& ctx) override;
    Result on_receive(context::Context& ctx, void *ptr, uint32_t sz,
                      uint32_t& sent) override;
    Result on_shutdown(context::Context& ctx) override;
     void on_delete(context::Context& ctx) override;

private:
    Connection *out_link;
};

} // namespace connection

} // namespace mesh

#endif // MESH_MULTIPOINT_H
