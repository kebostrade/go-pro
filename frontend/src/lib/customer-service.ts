'use client';

import {
  collection,
  doc,
  addDoc,
  updateDoc,
  deleteDoc,
  getDoc,
  getDocs,
  query,
  orderBy,
  where,
  serverTimestamp,
} from 'firebase/firestore';
import { getDbInstance } from '@/lib/firebase';
import { Customer, CustomerFormData } from '@/types/customer';

const CUSTOMERS_COLLECTION = 'customers';

export const customerService = {
  // Get all customers
  async getAll(): Promise<Customer[]> {
    const db = getDbInstance();
    const q = query(
      collection(db, CUSTOMERS_COLLECTION),
      orderBy('createdAt', 'desc')
    );
    
    const snapshot = await getDocs(q);
    return snapshot.docs.map((doc) => ({
      id: doc.id,
      ...doc.data(),
    })) as Customer[];
  },

  // Get customer by ID
  async getById(id: string): Promise<Customer | null> {
    const db = getDbInstance();
    const docRef = doc(db, CUSTOMERS_COLLECTION, id);
    const snapshot = await getDoc(docRef);
    
    if (!snapshot.exists()) {
      return null;
    }
    
    return { id: snapshot.id, ...snapshot.data() } as Customer;
  },

  // Create new customer
  async create(data: CustomerFormData): Promise<Customer> {
    const db = getDbInstance();
    const now = new Date().toISOString();
    
    const customerData = {
      ...data,
      totalOrders: 0,
      totalSpent: 0,
      createdAt: now,
      updatedAt: now,
    };
    
    const docRef = await addDoc(collection(db, CUSTOMERS_COLLECTION), customerData);
    
    return {
      id: docRef.id,
      ...customerData,
    };
  },

  // Update customer
  async update(id: string, data: Partial<CustomerFormData>): Promise<void> {
    const db = getDbInstance();
    const docRef = doc(db, CUSTOMERS_COLLECTION, id);
    
    await updateDoc(docRef, {
      ...data,
      updatedAt: new Date().toISOString(),
    });
  },

  // Delete customer
  async delete(id: string): Promise<void> {
    const db = getDbInstance();
    const docRef = doc(db, CUSTOMERS_COLLECTION, id);
    await deleteDoc(docRef);
  },

  // Search customers
  async search(searchTerm: string): Promise<Customer[]> {
    const allCustomers = await this.getAll();
    const term = searchTerm.toLowerCase();
    
    return allCustomers.filter((customer) =>
      customer.firstName?.toLowerCase().includes(term) ||
      customer.lastName?.toLowerCase().includes(term) ||
      customer.email?.toLowerCase().includes(term) ||
      customer.company?.toLowerCase().includes(term) ||
      customer.phone?.includes(term)
    );
  },

  // Filter customers by status
  async getByStatus(status: 'active' | 'inactive' | 'prospect'): Promise<Customer[]> {
    const db = getDbInstance();
    const q = query(
      collection(db, CUSTOMERS_COLLECTION),
      where('status', '==', status),
      orderBy('createdAt', 'desc')
    );
    
    const snapshot = await getDocs(q);
    return snapshot.docs.map((doc) => ({
      id: doc.id,
      ...doc.data(),
    })) as Customer[];
  },
};
