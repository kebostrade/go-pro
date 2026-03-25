// Customer types for the management system
export interface Customer {
  id: string;
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  company?: string;
  address?: string;
  city?: string;
  country?: string;
  postalCode?: string;
  notes?: string;
  status: 'active' | 'inactive' | 'prospect';
  tags?: string[];
  totalOrders?: number;
  totalSpent?: number;
  createdAt: string;
  updatedAt: string;
}

export interface CustomerFormData {
  firstName: string;
  lastName: string;
  email: string;
  phone: string;
  company: string;
  address: string;
  city: string;
  country: string;
  postalCode: string;
  notes: string;
  status: 'active' | 'inactive' | 'prospect';
  tags: string[];
}

export const initialCustomerFormData: CustomerFormData = {
  firstName: '',
  lastName: '',
  email: '',
  phone: '',
  company: '',
  address: '',
  city: '',
  country: '',
  postalCode: '',
  notes: '',
  status: 'prospect',
  tags: [],
};
